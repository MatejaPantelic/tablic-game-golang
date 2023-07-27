package main

import (
	"github.com/gin-gonic/gin"
	"main.go/api"
	"main.go/initializers"
	"main.go/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	r := gin.Default()

	r.POST("/addPlayer", api.AddPlayerHandler)
	r.GET("/cards", api.NewDeckHandler)
	r.GET("/cards/:userId/:deckId", middleware.CheckAuthTokenUserId, api.ShowPlayerCards)
	r.GET("/takecardsfromtable/:deckId/:handPile/:takenPile", middleware.CheckAuthTokenDeckId, api.TakeCardsFromTable)
	r.GET("/throwCard/:cardCode/:deckId/:playerPile", middleware.CheckAuthTokenDeckId, api.ThrowCardHandler)
	r.GET("/gettoken/:userId/:deckId", api.MakeToken)

	r.Run()

}