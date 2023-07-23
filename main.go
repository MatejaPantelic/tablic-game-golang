package main

import (
	"github.com/gin-gonic/gin"
	"main.go/api"
	"main.go/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	r := gin.Default()
	api.InitializeHandlers(r)

	r.Run(":8080")

}
