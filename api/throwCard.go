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


	func throwCardHandler(c *gin.Context){
	cardCode :=c.Param("cardCode")
	deckId :=c.Param("deckId")
	playerPile :=c.Param("playerPile")

	var exist bool= existsInDeck(cardCode) //checking if card exist in deck
	if exist{
		playerCards, _ := http.Get(fmt.Sprintf(constants.ListPileCardsURL, deckId, playerPile))	//list of cards in a player hand pile
		body := parseJsonToStruct(playerCards)
		var cardInPiles models.ListCardResponse
		err :=json.Unmarshal(body, &cardInPiles)
		if(err != nil){
			log.Fatal(err)
		}
		existInHand:=IsCardExistInHand(playerPile,cardInPiles,cardCode)

		if existInHand {
			http.Get(fmt.Sprintf(constants.DrawCardsFromPileURL, deckId, playerPile,cardCode)) //drawing from pile
			http.Get(fmt.Sprintf(constants.AddToPileUrl, deckId,"table",cardCode)) //adding card to table pile
			existInHand=false
			c.JSON(http.StatusOK, gin.H{"response": "The card is thrown on the table"})

		}else{
			c.JSON(http.StatusForbidden, gin.H{"response": "The selected card is not in your hand."})
		}

	}else{
		c.JSON(http.StatusForbidden, gin.H{"response": "The selected card does not exist in the deck."})
	}
		
}	//Function for checking if cards exists in the player's hand
	func IsCardExistInHand(playerPile string, cardInPiles models.ListCardResponse, cardCode string) bool {
		if (playerPile=="hand1"){		
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
