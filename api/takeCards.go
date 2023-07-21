package api

import (
	"github.com/gin-gonic/gin"
	"main.go/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"fmt"
	"strings"
	"strconv"
	"main.go/utilities"
)

//Function for listing cards in a pile
func listPileCards(deck string, pileName string)(CardsArray []models.CardList){
	url := fmt.Sprintf(utilities.LIST_PILE_CARDS_URL, deck, pileName)
	resp, errURL := http.Get(url)
	if errURL != nil {
		log.Fatal(errURL)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var ListCardResponse models.ListCardResponse
	json.Unmarshal(body, &ListCardResponse)

	switch(pileName){
	case "hand1": CardsArray = ListCardResponse.Piles.Hand1.Cards
	case "hand2": CardsArray = ListCardResponse.Piles.Hand2.Cards
	case "table": CardsArray = ListCardResponse.Piles.Table.Cards
	default: 
	}

	return
}

//Function for draeing cards from a pile
func drawCardsFromPile(deck string, pileName string, cards string){
	url := fmt.Sprintf(utilities.DRAW_CARDS_FROM_PILE_URL, deck, pileName, cards)
	resp, errURL := http.Get(url)
	if errURL != nil {
		log.Fatal(errURL)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var DrowCardResponse models.DrawingFromPilesResponse
	json.Unmarshal(body, &DrowCardResponse)
}

//Function for adding cards to a pile
func addToPile(deck string, pileName string, cards string){
	url := fmt.Sprintf(utilities.ADD_TO_PILE_URL, deck, pileName, cards)
	resp, errURL := http.Get(url)
	if errURL != nil {
		log.Fatal(errURL)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var AddCardsResponse models.AddingToPilesResponse
	json.Unmarshal(body, &AddCardsResponse)
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
	}
	HandCard := RequestData.HandCard
	TakenCardsString := RequestData.TakenCards
	TakenCards := strings.Split(TakenCardsString, ",")

	//VALIDATE CARD FROM HAND
	var HandCards []models.CardList = listPileCards(deckId, handPile)
	var existInHand bool = false
	for _, card := range HandCards {
		if card.Code == HandCard{
			existInHand=true
		}
	}
	if(!existInHand){
		c.JSON(http.StatusNotFound, gin.H{"response": "The selected card is not in your hand."})
	}

	//VALIDATE CARDS FROM TABLE
	var TableCards []models.CardList = listPileCards(deckId, "table")
	for _, cardTaken := range TakenCards {
		existInHand = false
		for _, cardTable := range TableCards{
			if cardTaken == cardTable.Code{
				existInHand = true
			}
		}
		if(existInHand == false)	{
			c.JSON(http.StatusNotFound, gin.H{"response": "Some of selected cards is not on the table."})
		}	
	}

	//CHECK VALUES
	var valid bool = false
	var HandCardValue int
	var sum int = 0
	switch(HandCard[0]){
	case 'J': HandCardValue = 12
	case 'Q': HandCardValue = 13
	case 'K': HandCardValue = 14
	default: HandCardValue,_ = strconv.Atoi(string(HandCard[0]))
	}
	fmt.Println(HandCardValue)
	for _, cardTaken := range TakenCards{
		var val int
		switch(cardTaken[0]){
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

		c.JSON(http.StatusOK, gin.H{"response": "Cards are moved from hand and table to taken pile"})
	}else{
		c.JSON(http.StatusNotFound, gin.H{"response": "You can't take chosen cards"})
	}

}