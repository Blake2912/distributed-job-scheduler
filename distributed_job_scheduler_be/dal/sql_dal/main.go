package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/sql_dal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	sql, err := config.DB.DB()

	if err != nil {
		log.Fatal(err)
	}
	// Ensures that we close the connection always when application shuts down
	defer config.CloseSqlConnection(sql)

	err = config.Migrate(config.DB)
	if err != nil {
		log.Fatal("Error occurred during migration", err)
	}
	log.Println("Migrations completed")
	
	router := gin.Default()
	router.GET("hello/", hello)
	router.Run("localhost:8081")
}

func hello(c *gin.Context){
	msg := fmt.Sprintln("Hello")
	c.IndentedJSON(http.StatusOK, msg)
}
