package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Pong main",
		})
	})

	fmt.Println("Server Starting on http://localhost:8080")
	router.Run(":8080")
}
