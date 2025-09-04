package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"

	"github.com/hiromaily/river-example/internal/jobs"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbURL := mustGetenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("pgxpool.New: %v", err)
	}
	defer pool.Close()

	// Worker 登録
	workers := river.NewWorkers()
	river.AddWorker(workers, &jobs.EmailWorker{})

	// River クライアント作成（default キューに最大 10 並列）
	rClient, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 10},
		},
		Workers: workers,
	})
	if err != nil {
		log.Fatalf("river.NewClient: %v", err)
	}

	// ワーカー開始
	go func() {
		if err := rClient.Start(ctx); err != nil {
			log.Fatalf("river client start: %v", err)
		}
	}()

	// HealthCheck認用 HTTP（任意）
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	go func() {
		log.Printf("worker http listening on :8081")
		if err := http.ListenAndServe(":8081", nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http: %v", err)
		}
	}()

	// SIGTERM/SIGINT でグレースフル停止
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
	log.Println("shutting down ...")
	shutdownCtx, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()
	if err := rClient.Stop(shutdownCtx); err != nil { // 進行中ジョブを待って停止
		log.Printf("river stop: %v", err)
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("required env %s is empty", k)
	}
	return v
}
