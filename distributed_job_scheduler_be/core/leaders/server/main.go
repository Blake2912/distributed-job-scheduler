package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/leader/internal/api/routes"
	"example.com/leader/internal/container"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	fmt.Println("Hello from leader")

	// Build required classes for application startup
	container := container.BuildContainer()

	// Register required routes
	routes.RegisterRoutes(router, container)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error:", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received")

	shutdownCtx, shutDownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutDownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("Server forced to shutdown:", err)
	}

	log.Println("Server exited cleanly")

}
