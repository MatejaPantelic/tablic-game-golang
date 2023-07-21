package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/initializers"
	"main.go/models"
)

func showPlayerCards(c *gin.Context) {
	//function recieve two parameter from endpoint for list card for player and table----USER ID AND DECK ID
	//we use Deck ID for recognition game of player
	userid := c.Param("userid")
	deckid := c.Param("deckid")

	// create variable type of structure Game to store data from database
	var game models.Game

	//variable for handling http request error
	var err error

	//take information for endpoint : deck which user use and name of hand pile for retrieve information from external api
	initializers.DB.Where("user_id = ? AND deck_pile = ?", userid, deckid).Find(&game)

	//call endpoint for list hand cards with necessary information DECK ID and NAME OF HAND PILE used from variable game
	respHand, err := http.Get("https://www.deckofcardsapi.com/api/deck/" + game.DeckPile + "/pile/" + game.HandPile + "/list/")

	//handle if there some error from nttp
	if err != nil {
		log.Fatal(err)
	}

	//part of code for accepting response
	defer respHand.Body.Close()
	bodyHand, _ := io.ReadAll(respHand.Body)

	//variable for store reponse from http request in acceptable format
	var drawResponse models.DrawResponse
	json.Unmarshal(bodyHand, &drawResponse)

	//variable to store cards for only requested player.
	var handcardsarray []models.CardList

	//http reponse return both json object(hand1 & hand2), we looks for one
	if drawResponse.Piles.Hand1.Cards == nil {
		handcardsarray = drawResponse.Piles.Hand2.Cards
	} else if drawResponse.Piles.Hand2.Cards == nil {
		handcardsarray = drawResponse.Piles.Hand1.Cards
	}

	//call endpoint for list table cards with necessary information DECK ID and NAME OF TABLE PILE used from variable game
	respDeck, err := http.Get("https://www.deckofcardsapi.com/api/deck/" + game.DeckPile + "/pile/" + game.TablePile + "/list/")

	//handle if there some error from nttp
	if err != nil {
		log.Fatal(err)
	}

	//part of code for accepting response
	defer respDeck.Body.Close()
	bodyDeck, _ := io.ReadAll(respDeck.Body)

	//variable for store reponse from http request in acceptable format
	var drawResponseDeck models.DrawResponse
	json.Unmarshal(bodyDeck, &drawResponseDeck)

	//return response with needed information
	c.JSON(http.StatusOK, gin.H{"User hand cards": handcardsarray, "Cards from table": drawResponseDeck.Piles.Table})

}

func InitializeHendlers(router *gin.Engine) {
	router.GET("/cards/:userid/:deckid", showPlayerCards)
}
