package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/proxy-service/internal/handler"
	"github.com/Blake2912/distributed-job-scheduler/proxy-service/internal/resolver"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/redis_dal/commands/queries"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/redis_dal/redisclient"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	redisCtx, rediscancel := context.WithTimeout(ctx, 20*time.Second)
	defer rediscancel()

	rdb, err := redisclient.New(redisCtx)
	if err != nil {
		log.Fatal(err)
	}
	defer rdb.Close()
	log.Println("Connected to redis")

	redisQueries := queries.NewRedisQueries(rdb)
	leaderResolver := resolver.NewLeaderResovler(redisQueries)
	proxyHandler := handler.NewProxyHandler(leaderResolver)

	mux := http.NewServeMux()
	mux.HandleFunc("/", proxyHandler.ProxyHandler())

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down proxy server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
