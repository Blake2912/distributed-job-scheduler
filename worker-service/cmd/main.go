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

	"github.com/Blake2912/distributed-job-scheduler/worker-service/internal/container"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	container := container.BuildContainer(
		httpClient,
		"http://:8081", //proxy/scheduler URL
		3,              // no. of workers
	)

	log.Println("Worker service started")

	container.App.Run(ctx)

	log.Println("Worker service stopped")
	/*
		router := gin.Default()
		router.GET("/hello", hello)
		router.Run(":8082")
	*/
}

func hello(c *gin.Context) {
	podName := os.Getenv("POD_NAME")
	podId := os.Getenv("POD_UID")
	msg := fmt.Sprintf("Hello from worker | pod id: %s | pod name: %s", podId, podName)
	c.IndentedJSON(http.StatusOK, msg)
}
