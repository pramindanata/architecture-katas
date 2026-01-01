package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

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

	err = rdb.Set(ctx, "ticket:a", 10_000_000, 0).Err()

	if err != nil {
		log.Fatalf("Failed to set ticket:a in Redis: %v", err)
	}

	err = rdb.Set(ctx, "ticket:b", 10_000_000, 0).Err()

	if err != nil {
		log.Fatalf("Failed to set ticket:b in Redis: %v", err)
	}

	fmt.Println("Initial values set in Redis")
}
