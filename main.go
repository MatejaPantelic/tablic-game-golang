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

	r.POST("/addPlayer", api.AddPlayerHandler)
	r.GET("/cards", api.NewDeckHandler)
	r.GET("/cards/:userId/:deckId", authorization.CheckAuthTokenUser, api.ShowPlayerCards)
	// r.GET("/takecardsfromtable/:deckId/:handPile/:takenPile", authorization.CheckAuthTokenUser, api.TakeCardsFromTable)
	r.GET("/takecardsfromtable/:deckId/:handPile/:takenPile", api.TakeCardsFromTable)
	// r.GET("/throwCard/:cardCode/:deckId/:playerPile", authorization.CheckAuthTokenUser, api.ThrowCardHandler)
	r.GET("/throwCard/:cardCode/:deckId/:playerPile", api.ThrowCardHandler)
	r.GET("/gettoken/:userId/:deckId", authorization.MakeToken)

	r.Run()

}
