package main

import (
	"github.com/gin-gonic/gin"
	"main.go/api"
	"main.go/authorization"
	"main.go/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	r := gin.Default()

	// Middleware applied to all routes within the /cards group
	// cardsGroup := r.Group("")
	// cardsGroup.Use(authorization.CheckAuthTokenUser("deckId"))
	// cardsGroup.GET("/throwCard/:cardCode/:deckId/:playerPile", api.ShowPlayerCards)

	r.POST("/addPlayer", api.AddPlayerHandler)
	r.GET("/cards", api.NewDeckHandler)
	r.GET("/cards/:userId/:deckId", api.ShowPlayerCards)
	r.GET("/takecardsfromtable/:deckId/:handPile/:takenPile", api.TakeCardsFromTable)
	r.GET("/throwCard/:cardCode/:deckId/:playerPile", api.ThrowCardHandler)
	r.GET("/gettoken/:userId/:deckId", authorization.MakeToken)

	r.Run()

}
