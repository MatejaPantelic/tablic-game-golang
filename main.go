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
  r.GET("/cards", api.newDeckHandler)
	r.GET("/cards/:userid/:deckid", api.ShowPlayerCards)
	r.GET("/takecardsfromtable/:deckId/:handPile/:takenPile", api.TakeCardsFromTable)
	r.GET("/throwCard/:cardCode/:deckId/:playerPile", api.throwCardHandler)
	r.Run()
}
