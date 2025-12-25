package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/hello", hello)
	router.Run(":8080")
}

func hello(c *gin.Context) {
	podName := os.Getenv("POD_NAME")
	podId := os.Getenv("POD_UID")
	msg := fmt.Sprintf("Hello from worker | pod id: %s | pod name: %s", podId, podName)
	c.IndentedJSON(http.StatusOK, msg)
}
