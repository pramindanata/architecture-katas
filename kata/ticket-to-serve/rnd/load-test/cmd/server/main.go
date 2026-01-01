package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

type Ticket struct {
	ID    string `json:"id"`
	Count int    `json:"count"`
}

type RequestBody struct {
	Tickets []Ticket `json:"tickets"`
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbName := os.Getenv("PG_DBNAME")
	dbUser := os.Getenv("PG_USER")
	dbPassword := os.Getenv("PG_PASSWORD")
	dbHost := "pg"
	dbPort := os.Getenv("PG_PORT")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
	)

	config, err := pgxpool.ParseConfig(psqlInfo)

	if err != nil {
		log.Fatalf("Unable to parse config: %v\n", err)
	}

	config.MaxConns = 50
	config.MinConns = 25
	config.MaxConnIdleTime = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	defer pool.Close()

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDBStr := os.Getenv("REDIS_DB")
	redisDB, err := strconv.Atoi(redisDBStr)

	if err != nil {
		log.Fatalf("Invalid REDIS_DB value: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       redisDB,
	})

	mux := http.NewServeMux()

	mux.HandleFunc("POST /pg-reserve-ticket", func(w http.ResponseWriter, r *http.Request) {
		var body RequestBody
		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %s", err), http.StatusBadRequest)
			return
		}

		if len(body.Tickets) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Bad request: no tickets provided")
			return
		}

		tx, err := pool.Begin(r.Context())

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to begin transaction: %s", err), http.StatusInternalServerError)
			return
		}

		defer tx.Rollback(r.Context())

		// Short term solution to prevent PG to do sequential scans.
		// When the available data is less than 700K, PG tends to do sequential scans.
		// It caused the p95 latency to spike from 15ms to 2s.
		_, err = tx.Exec(r.Context(), "SET LOCAL enable_seqscan = off;")

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to set enable_seqscan: %s", err), http.StatusInternalServerError)
			return
		}

		for _, ticket := range body.Tickets {
			var ticketUnitIDs []int
			query := `
				SELECT id FROM ticket_units
				WHERE ticket_id = $1 AND reserved = false
				LIMIT $2
				FOR UPDATE SKIP LOCKED
			`
			rows, err := tx.Query(r.Context(), query, ticket.ID, ticket.Count)

			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to query ticket units: %s", err), http.StatusInternalServerError)
				return
			}

			for rows.Next() {
				var id int

				if err := rows.Scan(&id); err != nil {
					rows.Close()
					http.Error(w, fmt.Sprintf("Failed to scan ticket unit id: %s", err), http.StatusInternalServerError)
					return
				}

				ticketUnitIDs = append(ticketUnitIDs, id)
			}

			// Check for errors during iteration
			if err := rows.Err(); err != nil {
				rows.Close()
				http.Error(w, fmt.Sprintf("Error iterating rows: %s", err), http.StatusInternalServerError)
				return
			}

			// CRITICAL: Close rows immediately after use, not deferred
			rows.Close()

			if len(ticketUnitIDs) < ticket.Count {
				http.Error(w, fmt.Sprintf("Not enough stock for ticket %s", ticket.ID), http.StatusConflict)
				return
			}

			updateQuery := "UPDATE ticket_units SET reserved = true WHERE id = ANY($1)"
			_, err = tx.Exec(r.Context(), updateQuery, ticketUnitIDs)

			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to reserve tickets: %s", err), http.StatusInternalServerError)
				return
			}
		}

		if err := tx.Commit(r.Context()); err != nil {
			http.Error(w, fmt.Sprintf("Failed to commit transaction: %s", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Tickets reserved successfully")
	})

	mux.HandleFunc("POST /complex-redis-decrement", func(w http.ResponseWriter, r *http.Request) {
		var body RequestBody
		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %s", err), http.StatusBadRequest)
			return
		}

		if len(body.Tickets) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Bad request: no tickets provided")
			return
		}

		keys := make([]string, len(body.Tickets))
		args := make([]interface{}, len(body.Tickets))

		for i, ticket := range body.Tickets {
			keys[i] = "ticket:" + ticket.ID
			args[i] = ticket.Count
		}

		luaScript := `
			for i=1, #KEYS do
				local current = redis.call('GET', KEYS[i])
				if current == false then
					return {err = "Ticket " .. KEYS[i] .. " does not exist"}
				end
				local current_val = tonumber(current)
				local decrement_val = tonumber(ARGV[i])
				if current_val < decrement_val then
					return {err = "Not enough stock for ticket " .. KEYS[i]}
				end
			end

			local results = {}
			for i=1, #KEYS do
				local decrement_val = tonumber(ARGV[i])
				results[i] = redis.call('DECRBY', KEYS[i], decrement_val)
			end

			return results
		`

		script := redis.NewScript(luaScript)
		_, err = script.Run(ctx, rdb, keys, args...).Result()

		if err != nil {
			message := fmt.Sprintf("Failed to decrement tickets: %v", err)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "Success")
	})

	mux.HandleFunc("POST /simple-redis-decrement", func(w http.ResponseWriter, r *http.Request) {
		var body RequestBody
		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %s", err), http.StatusBadRequest)
			return
		}

		if len(body.Tickets) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Bad request: no tickets provided")
			return
		}

		keys := make([]string, len(body.Tickets))
		args := make([]interface{}, len(body.Tickets))

		for i, ticket := range body.Tickets {
			keys[i] = "ticket:" + ticket.ID
			args[i] = ticket.Count
		}

		luaScript := `
			local results = {}
			for i=1, #KEYS do
				local decrement_val = tonumber(ARGV[i])
				results[i] = redis.call('DECRBY', KEYS[i], decrement_val)
			end
			return results
		`
		script := redis.NewScript(luaScript)
		_, err = script.Run(ctx, rdb, keys, args...).Result()

		if err != nil {
			message := fmt.Sprintf("Failed to decrement tickets: %v", err)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "Success")
	})

	serverPort := os.Getenv("SERVER_PORT")
	fmt.Printf("Server starting on port %s\n", serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", serverPort), mux))
}
