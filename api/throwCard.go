package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/constants"
	"main.go/initializers"
	"main.go/models"
)

func ThrowCardHandler(c *gin.Context) {
	cardCode := c.Param("cardCode")
	deckId := c.Param("deckId")
	playerPile := c.Param("playerPile")
	
	// create variable type of structure Game
	var game models.Game
	result := initializers.DB.Model(&game).Where("hand_pile = ? AND deck_pile = ?", playerPile, deckId).Find(&game)
	errorCheck(result.Error, 400, "Failed to fetch data from DB",c)
	//checking if it is the player's turn to play
	if game.First{
		var exist bool = existsInDeck(cardCode) 
		//checking if card exist in deck
		if exist {
			//list of cards in a player hand pile
			cardInPiles := getCardsFromPile(deckId,playerPile)

			existInHand := isCardInHand(playerPile, cardInPiles, cardCode)
			//checking if the card is in the player's hand
			if existInHand {
				//adding card to table pile
				http.Get(fmt.Sprintf(constants.ADD_TO_PILE_URL, deckId, "table", cardCode)) 
				existInHand = false
				c.JSON(http.StatusOK, gin.H{
					"message": "The card is thrown on the table", 
					"user_hand_cards": getCardsFromPile(deckId,playerPile).Piles,
					"table_cards": getCardsFromPile(deckId,"table").Piles.Table,
				})
				// create variable type of structure Game
				whoPlaysNext(c, playerPile, deckId)

				FinishGame(c, deckId)

			} else {
				c.JSON(http.StatusOK, gin.H{"response": "The selected card is not in your hand."})
			}

		} else {
			c.JSON(http.StatusForbidden, gin.H{"response": "The selected card does not exist in the deck."})
	}
	}else{
		c.JSON(http.StatusBadRequest, gin.H{"response": "The opponent play next."})
	}
} 

//Function for checking if cards exists in the player's hand
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
