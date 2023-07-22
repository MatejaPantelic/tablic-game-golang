package api

import (
	"github.com/gin-gonic/gin"
	"main.go/models"
	"encoding/json"
	"log"
	"net/http"
	"fmt"
	"strings"
	"strconv"
	"main.go/constants"
)

//Function for listing cards in a pile
func listPileCards(deck string, pileName string)(CardsArray []models.CardList){
	url := fmt.Sprintf(constants.ListPileCardsURL, deck, pileName)
	resp, errURL := http.Get(url)
	if errURL != nil {
		log.Fatal(errURL)
	}

	body := parseJsonToStruct(resp) 

	var ListCardResponse models.ListCardResponse
	err := json.Unmarshal(body, &ListCardResponse)
	if(err != nil){
		log.Fatal(err)
	}

	switch(pileName){
	case "hand1": CardsArray = ListCardResponse.Piles.Hand1.Cards
	case "hand2": CardsArray = ListCardResponse.Piles.Hand2.Cards
	case "table": CardsArray = ListCardResponse.Piles.Table.Cards
	default: 
	}

	return
}

//Function for drawing cards from a pile
func drawCardsFromPile(deck string, pileName string, cards string){
	url := fmt.Sprintf(constants.DrawCardsFromPileURL, deck, pileName, cards)
	resp, errURL := http.Get(url)
	if errURL != nil {
		log.Fatal(errURL)
	}

	body := parseJsonToStruct(resp) 

	var DrowCardResponse models.DrawingFromPilesResponse
	err := json.Unmarshal(body, &DrowCardResponse)
	if(err != nil){
		log.Fatal(err)
	}
}

//Function for adding cards to a pile
func addToPile(deck string, pileName string, cards string){
	url := fmt.Sprintf(constants.AddToPileUrl, deck, pileName, cards)
	resp, errURL := http.Get(url)
	if errURL != nil {
		log.Fatal(errURL)
	}

	body := parseJsonToStruct(resp) 

	var AddCardsResponse models.AddingToPilesResponse
	err := json.Unmarshal(body, &AddCardsResponse)
	if(err != nil){
		log.Fatal(err)
	}

}

type RequestData struct{
	HandCard string `json:"hand_card"`
	TakenCards string `json:"taken_cards"`
}


func TakeCardsFromTable(c *gin.Context){
	//EXTRACT PARAMETERS
	deckId := c.Param("deckId")
	handPile := c.Param("handPile")
	takenPile := c.Param("takenPile")

	//EXTRACT BODY REQUEST
	var RequestData RequestData
	err := c.BindJSON(&RequestData)
	if(err != nil){
		c.JSON(http.StatusBadRequest, gin.H{"response": "Invalid JSON format in request body"})
		return
	}
	HandCard := strings.ToUpper(RequestData.HandCard)
	TakenCardsString := strings.ToUpper(RequestData.TakenCards)
	TakenCards := strings.Split(TakenCardsString, ",")

	//VALIDATE CARD FROM HAND
	if(!existsInDeck(HandCard)){
		c.JSON(http.StatusForbidden, gin.H{"response": "The selected hand card does not exist in the deck."})
		return
	}

	var HandCards []models.CardList = listPileCards(deckId, handPile)

	if(!existsInPile(HandCard, HandCards)){
		c.JSON(http.StatusForbidden, gin.H{"response": "The selected card is not in your hand."})
		return
	}

	//VALIDATE CARDS FROM TABLE
	var TableCards []models.CardList = listPileCards(deckId, "table")
	for _, cardTaken := range TakenCards {
		if(!existsInDeck(cardTaken)){
			c.JSON(http.StatusForbidden, gin.H{"response": "The selected table card does not exist in the deck."})
			return
		}
		
		if(!existsInPile(cardTaken, TableCards))	{
			c.JSON(http.StatusForbidden, gin.H{"response": "Some of selected cards is not on the table."})
			return
		}	
	}

	//CHECK VALUES
	var valid bool = false
	var HandCardValue int
	var sum int = 0
	switch(HandCard[0]){
	case '0': HandCardValue = 10
	case 'J': HandCardValue = 12
	case 'Q': HandCardValue = 13
	case 'K': HandCardValue = 14
	default: HandCardValue,_ = strconv.Atoi(string(HandCard[0]))
	}
	fmt.Println(HandCardValue)
	for _, cardTaken := range TakenCards{
		var val int
		switch(cardTaken[0]){
			case '0': val = 10
			case 'J': val = 12
			case 'Q': val = 13
			case 'K': val = 14
			default: val,_ = strconv.Atoi(string(cardTaken[0]))
		}
		sum += val
		fmt.Println(val, sum)
	}
	if (sum%HandCardValue == 0){
		valid = true
	}

	//IF VALID MOVE CARDS FROM HAND AND TABLE PILE TO TAKEN PILE
	if(valid){
		drawCardsFromPile(deckId, handPile, HandCard)
		separator := ","
		cards := strings.Join(TakenCards, separator)
		drawCardsFromPile(deckId, "table", cards)
		addToPile(deckId, takenPile, cards+","+HandCard)

		c.JSON(http.StatusOK, gin.H{"response": "Cards are moved from hand and table pile to taken pile"})
	}else{
		c.JSON(http.StatusNotFound, gin.H{"response": "You can't take chosen cards"})
	}

}