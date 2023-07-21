package main

import (
	// "encoding/json"
	// "fmt"
	// "io"
	// "net/http"

	"github.com/gin-gonic/gin"
	"main.go/api"
	"main.go/initializers"
	// "main.go/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

// func showPlayerCards(c *gin.Context) {
// 	userid := c.Param("userid")
// 	deckid := c.Param("deckid")
// 	var game models.Game
// 	// initializers.DB.First(&game, idStr)
// 	initializers.DB.Where("user_id = ? AND deck_pile = ?", userid, deckid).Find(&game)

// 	resp, _ := http.Get("https://www.deckofcardsapi.com/api/deck/" + game.DeckPile + "/pile/" + game.HandPile + "/list/")

// 	defer resp.Body.Close()
// 	body, _ := io.ReadAll(resp.Body)
// 	var drawResponse models.DrawResponse
// 	json.Unmarshal(body, &drawResponse)
// 	var handcardsarray []models.CardList
// 	if drawResponse.Piles.Hand1.Cards == nil {
// 		fmt.Println("Function1 is nil")
// 		handcardsarray = drawResponse.Piles.Hand2.Cards
// 	} else if drawResponse.Piles.Hand2.Cards == nil {
// 		fmt.Println("Function2 is nil")
// 		handcardsarray = drawResponse.Piles.Hand1.Cards
// 	}

// 	resp2, _ := http.Get("https://www.deckofcardsapi.com/api/deck/" + game.DeckPile + "/pile/" + game.TablePile + "/list/")
// 	defer resp2.Body.Close()
// 	body2, _ := io.ReadAll(resp2.Body)
// 	var drawResponse2 models.DrawResponse
// 	json.Unmarshal(body2, &drawResponse2)
// 	var tablecardsarray = drawResponse2.Piles.Table

// 	c.JSON(http.StatusOK, gin.H{"tvoje karte": handcardsarray, "tretne na tabli": tablecardsarray})
// 	// c.JSON(200, gin.H{"posts": game.DeckPile})

// }

func main() {
	r := gin.Default()

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"messege": "pong",
	// 	})
	// })
	// r.GET("/cards/:userid/:deckid", showPlayerCards)
	api.InitializeHendlers(r)

	r.Run()
}
