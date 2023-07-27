
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"main.go/initializers"
	"main.go/constants"
	"main.go/models"
	"github.com/gin-gonic/gin"
)

//Function for parsing response JSON to Struct
func parseJsonToStruct(resp *http.Response)(body []byte){
	defer resp.Body.Close()

	body, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		log.Fatal(errBody)
	}

	return
}

//Function for checking if card exists in deck
func existsInDeck(cardCode string)(exist bool){
	var allCards = [52]string{"AS","2S","3S","4S","5S","6S","7S","8S","9S","0S","JS","QS","KS",
	"AD","2D","3D","4D","5D","6D","7D","8D","9D","0D","JD","QD","KD",
	"AC","2C","3C","4C","5C","6C","7C","8C","9C","0C","JC","QC","KC",
	"AH","2H","3H","4H","5H","6H","7H","8H","9H","0H","JH","QH","KH"}

	exist = false
	for i:=0; i < 52; i++ {
		if allCards[i]==cardCode{
			exist=true
		}
	}

	return
}

//Function for checking if card exists in specific pile
func existsInPile(cardCode string, pile []models.CardList)(exist bool){
	exist = false
	for _, card := range pile {
		if card.Code == cardCode{
			exist=true
		}
	}
	return
}

func getCardsFromPile(deckId string, playerPile string)(cardInPiles models.ListCardResponse){
	playerCards, _ := http.Get(fmt.Sprintf(constants.LIST_PILE_CARDS_URL, deckId, playerPile))			 
	body := parseJsonToStruct(playerCards)
	err := json.Unmarshal(body, &cardInPiles)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func whoPlaysNext(c *gin.Context, playerPile string, deckId string){
	var game models.Game
	//set attribute "first" on false for player 1
	result :=initializers.DB.Model(&game).Where("hand_pile = ? AND deck_pile = ?", playerPile, deckId).Update("first", false)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"message": result.Error})
	}
	//set attribute "first" on true for player 2
	result =initializers.DB.Model(&game).Where("hand_pile NOT IN (?) AND deck_pile = ?", playerPile, deckId).Update("first", true)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": result.Error})
	}
}

func Piles(deckId string, pile string) models.PileList {

	var listCardResponse = listCardResponseFunction(deckId, pile)
	var Pile models.PileList

	if pile == "taken1" {
		Pile = listCardResponse.Piles.Taken1
	} else if pile == "taken2" {
		Pile = listCardResponse.Piles.Taken2
	} else if pile == "player1" {
		Pile = listCardResponse.Piles.Hand1
	} else if pile == "player2" {
		Pile = listCardResponse.Piles.Hand2
	} else if pile == "table" {
		Pile = listCardResponse.Piles.Table
	}
	return Pile
}

func notEmptyHands(deckId string) bool {

	var Pile1 models.PileList = Piles(deckId, "hand1")
	remaining1 := Pile1.Remaining

	var Pile2 models.PileList = Piles(deckId, "hand2")
	remaining2 := Pile2.Remaining

	if remaining1 != 0 && remaining2 != 0 {
		return true
	}
	return false
}

func listCardResponseFunction(deckId string, pile string) models.ListCardResponse {
	url := fmt.Sprintf(constants.LIST_PILE_CARDS_URL, deckId, pile)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Issue with specified URL")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error during reading response!")
	}

	var listCardResponse models.ListCardResponse
	err = json.Unmarshal(body, &listCardResponse)

	if err != nil {
		log.Fatal("Failed to read body")
	}

	return listCardResponse
}

func emptyDeck(deckId string, pile string) bool {
	var listCardResponse = listCardResponseFunction(deckId, pile)
	remaining := listCardResponse.Remaining

	if remaining == 0 {
		return true
	}
	return false
}

func emptyTable(deckId string) bool {
	TablePile := Piles(deckId, "table")
	remainingTable := TablePile.Remaining

	if remainingTable == 0 {
		return true
	}
	return false
}

func calculateCardPoints(code string) int {

	if code[0] == 'A' || code[0] == 'K' || code[0] == 'Q' || code[0] == 'J' {
		return 1
	} else if code[0] == '2' && code[1] == 'C' {
		return 1
	} else if code[0] == '0' {
		if code[1] == 'D' {
			return 2
		}
		return 1
	}
	return 0
}