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
	r.GET("/cards", api.NewDeckHandler)
	r.GET("/cards/:userid/:deckid", api.ShowPlayerCards)
	r.GET("/takecardsfromtable/:deckId/:handPile/:takenPile", api.TakeCardsFromTable)
	r.GET("/throwCard/:cardCode/:deckId/:playerPile", api.ThrowCardHandler)
	r.Run()

}
