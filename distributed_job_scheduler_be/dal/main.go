package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/dal/redis_dal/redisclient"
	"example.com/dal/sql_dal/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	//connect Db
	config.ConnectDB()

	sql, err := config.DB.DB()

	if err != nil {
		log.Fatal(err)
	}
	// Ensures that we close the connection always when application shuts down
	defer config.CloseSqlConnection(sql)

	//run migrations
	err = config.Migrate(config.DB)
	if err != nil {
		log.Fatal("Error occurred during migration", err)
	}
	log.Println("Migrations completed")

	//connect redis
	rdb, err := redisclient.New()
	if err != nil {
		log.Fatal(err)
	}
	defer rdb.Close()
	log.Println("Connected to redis")

	//router
	router := gin.Default()
	router.GET("hello/", hello)
	router.GET("helloRedis/", redisTest(rdb))
	router.Run("localhost:8081")
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
