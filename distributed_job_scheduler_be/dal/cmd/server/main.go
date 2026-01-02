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

	"example.com/dal/internal/api/routes"
	"example.com/dal/internal/container"
	"example.com/dal/redis_dal/redisclient"
	"example.com/dal/sql_dal/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

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

	// Application startup
	container := container.BuildContainer(config.DB)

	//router
	router := gin.Default()
	router.GET("hello/", hello)
	router.GET("helloRedis/", redisTest(rdb))

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

func hello(c *gin.Context) {
	msg := fmt.Sprintln("Hello")
	c.IndentedJSON(http.StatusOK, msg)
}

// temp test
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
