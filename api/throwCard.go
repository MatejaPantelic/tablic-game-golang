package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/constants"
	"main.go/models"
)

func ThrowCardHandler(c *gin.Context) {
	cardCode := c.Param("cardCode")
	deckId := c.Param("deckId")
	playerPile := c.Param("playerPile")

	var exist bool = existsInDeck(cardCode) //checking if card exist in deck
	if exist {
		playerCards, _ := http.Get(fmt.Sprintf(constants.LIST_PILE_CARDS_URL, deckId, playerPile)) //list of cards in a player hand pile
		body := parseJsonToStruct(playerCards)
		var cardInPiles models.ListCardResponse
		err := json.Unmarshal(body, &cardInPiles)
		if err != nil {
			log.Fatal(err)
		}
		existInHand := isCardInHand(playerPile, cardInPiles, cardCode)

		if existInHand {
			http.Get(fmt.Sprintf(constants.ADD_TO_PILE_URL, deckId, "table", cardCode)) //adding card to table pile
			existInHand = false
			c.JSON(http.StatusOK, gin.H{"response": "The card is thrown on the table"})

		} else {
			c.JSON(http.StatusForbidden, gin.H{"response": "The selected card is not in your hand."})
		}

	} else {
		c.JSON(http.StatusForbidden, gin.H{"response": "The selected card does not exist in the deck."})
	}

} //Function for checking if cards exists in the player's hand
func isCardInHand(playerPile string, cardInPiles models.ListCardResponse, cardCode string) bool {
	if playerPile == "hand1" {
		size := len(cardInPiles.Piles.Hand1.Cards)
		for i := 0; i < size; i++ {
			if cardInPiles.Piles.Hand1.Cards[i].Code == cardCode {
				return true
			}
		}
	} else if playerPile == "hand2" {
		size := len(cardInPiles.Piles.Hand2.Cards)
		for i := 0; i < size; i++ {
			if cardInPiles.Piles.Hand2.Cards[i].Code == cardCode {
				return true
			}
		}
	}
	return false
}
