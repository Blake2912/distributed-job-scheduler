package main

import (
	"fmt"
    "github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	fmt.Println("Hello from leader")
	router.Run("localhost:8080")
}
