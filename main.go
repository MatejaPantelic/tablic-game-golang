package main

import (
	"github.com/gin-gonic/gin"
	"main.go/initializers"
)

func init() {
	initializers.LoadEnvVariables()

}

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"messege": "pong",
		})
	})

	r.Run()
}
