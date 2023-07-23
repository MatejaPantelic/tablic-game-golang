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
	r.GET("/cards/:userid/:deckid", api.ShowPlayerCards)
	r.Run()
}
