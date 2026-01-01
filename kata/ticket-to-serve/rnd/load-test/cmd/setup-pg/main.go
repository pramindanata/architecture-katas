package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbName := os.Getenv("PG_DBNAME")
	dbUser := os.Getenv("PG_USER")
	dbPassword := os.Getenv("PG_PASSWORD")
	dbHost := os.Getenv("PG_HOST")
	dbPort := os.Getenv("PG_PORT")

	// First, connect to postgres database to create target database if not exists
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
	)

	config, err := pgxpool.ParseConfig(psqlInfo)

	if err != nil {
		log.Fatalf("Unable to parse config: %v\n", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	log.Println("Creating database if not exists...")
	_, err = pool.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE \"%s\";", dbName))

	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			log.Fatalf("Unable to create database: %v\n", err)
		}

		log.Println("Database already exists, continuing...")
	} else {
		log.Println("Database created successfully.")
	}

	pool.Close()

	// Now connect to the target database
	psqlInfo = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
	)

	config, err = pgxpool.ParseConfig(psqlInfo)

	if err != nil {
		log.Fatalf("Unable to parse config: %v\n", err)
	}

	pool, err = pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	defer pool.Close()

	log.Println("Dropping composite index if exists...")
	_, err = pool.Exec(context.Background(), "DROP INDEX IF EXISTS idx_ticket_units_on_ticket_id_and_reserved;")

	if err != nil {
		log.Fatalf("Unable to drop index: %v\n", err)
	}

	log.Println("Dropping ticket_units table if exists...")
	_, err = pool.Exec(context.Background(), "DROP TABLE IF EXISTS ticket_units;")

	if err != nil {
		log.Fatalf("Unable to drop ticket_units table: %v\n", err)
	}

	log.Println("Dropping tickets table if exists...")
	_, err = pool.Exec(context.Background(), "DROP TABLE IF EXISTS tickets;")

	if err != nil {
		log.Fatalf("Unable to drop tickets table: %v\n", err)
	}

	setupSQL := `
		CREATE TABLE tickets (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255),
			stock INT
		);

		INSERT INTO tickets (id, name, stock) VALUES
		('a', 'Ticket 1', 100000),
		('b', 'Ticket 2', 100000),
		('c', 'Ticket 3', 100000),
		('d', 'Ticket 4', 100000),
		('e', 'Ticket 5', 100000),
		('f', 'Ticket 6', 100000),
		('g', 'Ticket 7', 100000),
		('h', 'Ticket 8', 100000),
		('i', 'Ticket 9', 100000),
		('j', 'Ticket 10', 100000);
	`

	log.Println("Creating tickets table and inserting initial data...")
	_, err = pool.Exec(context.Background(), setupSQL)

	if err != nil {
		log.Fatalf("Unable to execute setup SQL: %v\n", err)
	}

	log.Println("Tickets table created and data inserted successfully.")

	createTicketUnitsTableSQL := `
		CREATE TABLE ticket_units (
			id SERIAL PRIMARY KEY,
			ticket_id VARCHAR(255) REFERENCES tickets(id),
			reserved BOOLEAN DEFAULT FALSE
		);
	`
	log.Println("Creating ticket_units table...")

	_, err = pool.Exec(context.Background(), createTicketUnitsTableSQL)

	if err != nil {
		log.Fatalf("Unable to create ticket_units table: %v\n", err)
	}

	log.Println("Creating index on ticket_units(ticket_id, reserved)...")

	_, err = pool.Exec(context.Background(), "CREATE INDEX idx_ticket_units_on_ticket_id_and_reserved ON ticket_units (ticket_id, reserved)")

	if err != nil {
		log.Fatalf("Unable to create index: %v\n", err)
	}

	log.Println("Index created successfully.")

	// Step 6: Create ticket_units records
	log.Println("Inserting ticket_units records...")
	var ticketIDs []string
	rows, err := pool.Query(context.Background(), "SELECT id FROM tickets")

	if err != nil {
		log.Fatalf("Unable to query ticket IDs: %v\n", err)
	}

	for rows.Next() {
		var id string

		if err := rows.Scan(&id); err != nil {
			log.Fatalf("Unable to scan ticket ID: %v\n", err)
		}

		ticketIDs = append(ticketIDs, id)
	}

	rows.Close()

	totalUnits := 100000
	batchSize := 10000

	for _, ticketID := range ticketIDs {
		log.Printf("Inserting units for ticket ID %s...\n", ticketID)
		startTime := time.Now()

		for i := 0; i < totalUnits; i += batchSize {
			var valueStrings []string
			var valueArgs []interface{}

			for j := 0; j < batchSize; j++ {
				valueStrings = append(valueStrings, "($"+strconv.Itoa(len(valueArgs)+1)+", $"+strconv.Itoa(len(valueArgs)+2)+")")
				valueArgs = append(valueArgs, ticketID, false)
			}

			stmt := fmt.Sprintf("INSERT INTO ticket_units (ticket_id, reserved) VALUES %s", strings.Join(valueStrings, ","))
			_, err := pool.Exec(context.Background(), stmt, valueArgs...)

			if err != nil {
				log.Fatalf("Unable to insert ticket_units batch: %v\n", err)
			}

			log.Printf("Inserted batch for ticket ID %s (%d/%d)\n", ticketID, i+batchSize, totalUnits)
		}

		log.Printf("Finished inserting units for ticket ID %s in %v\n", ticketID, time.Since(startTime))
	}

	log.Println("All ticket_units records created successfully.")
}
