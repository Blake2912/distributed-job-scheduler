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

	_ "github.com/Blake2912/distributed-job-scheduler/scheduler-service/docs"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/eventbus"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/httpclient"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/api/routes"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/container"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/pod_library/client"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/redis_dal/redisclient"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/redis_dal/redissubscriber"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Job Scheduler API
// @version         1.0
// @description     API for scheduler, image and worker spawning
// @host            localhost:8081
// @BasePath        /api
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	//connect Db
	dbCtx, dbCancel := context.WithTimeout(ctx, 60*time.Second)
	defer dbCancel()

	if err := config.ConnectDB(dbCtx); err != nil {
		log.Fatal(err)
	}
	defer config.CloseSqlConnection()

	//connect redis
	redisCtx, rediscancel := context.WithTimeout(ctx, 20*time.Second)
	defer rediscancel()

	rdb, err := redisclient.New(redisCtx)
	if err != nil {
		log.Fatal(err)
	}
	defer rdb.Close()
	log.Println("Connected to redis")

	httpClient := httpclient.New(120 * time.Minute)

	k8sClient, err := client.New()
	if err != nil {
		log.Fatal(err)
	}

	// Application startup
	// Future Improvmement: Move the parameters into a struct to make it readable if the parameters grows.
	container := container.BuildContainer(config.DB, rdb, ctx, httpClient, k8sClient)

	// Event bus
	bus := eventbus.NewEventBus[eventbus.TTLExpiredEvent](500)

	redissubscriber.PublishRedisKeyExpiryEvent(rdb, ctx, bus)
	// Future improvment: Use DI
	container.TTLExpiryConsumer.StartTTLExpiryExecution(ctx, bus)

	//router
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/api/hello", hello)
	router.GET("/api/helloRedis", redisTest(rdb))

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

	// Run leader election after the server has started and registered its routes
	err = container.LeaderElector.Run(ctx, func(leaderCtx context.Context) {
		container.Scheduler.Run(leaderCtx)
	})

	if err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
	log.Println("Shutdown signal received")

	shutdownCtx, shutDownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutDownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("Server forced to shutdown:", err)
	}

	log.Println("Server exited cleanly")
}

// Hello godoc
// @Summary Health check
// @Description Simple hello endpoint
// @Tags system
// @Produce json
// @Success 200 {string} string
// @Router /hello [get]
func hello(c *gin.Context) {
	msg := fmt.Sprintln("Hello")
	c.IndentedJSON(http.StatusOK, msg)
}

// RedisTest godoc
// @Summary Test Redis connection - TEMP
// @Description Writes and reads a test key from Redis
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /helloRedis [get]
func redisTest(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if err := rdb.Set(ctx, "ping", "pong", 10*time.Second).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		val, err := rdb.Get(ctx, "ping").Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"redis": val,
		})
	}
}
