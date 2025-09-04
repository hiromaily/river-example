package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"

	"github.com/hiromaily/river-example/internal/jobs"
)

func main() {
	ctx := context.Background()
	dbURL := mustGetenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("pgxpool.New: %v", err)
	}
	defer pool.Close()

	// Worker 登録（Insert-only クライアントでも必要）
	workers := river.NewWorkers()
	river.AddWorker(workers, &jobs.EmailWorker{})

	// Insert-only クライアント（Start しない）
	rClient, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		Workers: workers,
	})
	if err != nil {
		log.Fatalf("river.NewClient: %v", err)
	}

	// 単発投入（非トランザクション）
	_, err = rClient.Insert(ctx, jobs.EmailArgs{
		To:          "alice@example.com",
		Subject:     "Hello from River",
		Body:        "This email was queued via River!",
		RequestedAt: time.Now(),
	}, nil)
	if err != nil {
		log.Fatalf("insert: %v", err)
	}
	log.Println("queued: email.send -> alice@example.com")
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("required env %s is empty", k)
	}
	return v
}
