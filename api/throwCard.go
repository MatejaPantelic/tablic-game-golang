package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/constants"
	"main.go/models"
)


	func throwCardHandler(c *gin.Context){
	cardCode :=c.Param("cardCode")
	deckId :=c.Param("deckId")
	playerPile :=c.Param("playerPile")

	var allCards = [52]string{"AS","2S","3S","4S","5S","6S","7S","8S","9S","0S","JS","QS","KS",
	"AD","2D","3D","4D","5D","6D","7D","8D","9D","0D","JD","QD","KD",
	"AC","2C","3C","4C","5C","6C","7C","8C","9C","0C","JC","QC","KC",
	"AH","2H","3H","4H","5H","6H","7H","8H","9H","0H","JH","QH","KH"}  //codes of all the cards in the deck

	var exist bool= false
	for i:=0; i < 52; i++ {
		if allCards[i]==cardCode{	//checking if the card exist
			exist=true
		}
	}
	if exist{
		playerCards, _ := http.Get(fmt.Sprintf(constants.ListCardInPile, deckId, playerPile))	//list of cards in a player hand pile

		defer playerCards.Body.Close()
		body, _ := io.ReadAll(playerCards.Body)
		var cardInPiles models.ListCardResponse
		json.Unmarshal(body, &cardInPiles)

		existInHand:=IsCardExistInHand(playerPile,cardInPiles,cardCode)

		if existInHand {
			http.Get(fmt.Sprintf(constants.DrawingFromPile, deckId, playerPile,cardCode)) //drawing from pile
			http.Get(fmt.Sprintf(constants.AddingToPile, deckId,cardCode)) //adding card to table pile
			existInHand=false
			c.JSON(http.StatusOK, gin.H{"response": "The card is thrown on the table"})

		}else{
			c.JSON(http.StatusForbidden, gin.H{"response": "The selected card is not in your hand."})
		}

	}else{
		c.JSON(http.StatusForbidden, gin.H{"response": "The selected card does not exist in the deck."})
	}
		
}
	func IsCardExistInHand(playerPile string, cardInPiles models.ListCardResponse, cardCode string) bool {
		if (playerPile=="hand1"){		//checking if the selected card exists in the player's hand
			size:=len(cardInPiles.Piles.Hand1.Cards)	
			for i:=0; i<size; i++ {
				if cardInPiles.Piles.Hand1.Cards[i].Code == cardCode{
					return true
				}
			}
		}else if playerPile=="hand2"{		
			size:=len(cardInPiles.Piles.Hand2.Cards)		
			for i:=0; i<size; i++{
				if cardInPiles.Piles.Hand2.Cards[i].Code == cardCode{
					return true
				}
			}
		}
		return false
	}
