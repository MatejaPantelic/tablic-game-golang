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
	var countA int = 0
	switch(HandCard[0]){
	case '0': HandCardValue = 10
	case 'A': HandCardValue = 11 //We will always consider A from hand as 11
	case 'J': HandCardValue = 12
	case 'Q': HandCardValue = 13
	case 'K': HandCardValue = 14
	default: HandCardValue,_ = strconv.Atoi(string(HandCard[0]))
	}

	//Count number of As in Taken pile
	for _, cardTaken := range TakenCards{
		if(cardTaken[0] == 'A'){
			countA++;
		}
	}

	//depending on nomber of As in Taken pile we have 5 cases
	//Case1: countA=0
	//Case2: countA=1 -> 2 combinations: 1             sum=1
	//                                   11            sum=11 
	
	//Case3: countA=2 -> 3 combinations: 1  1          sum=2
	//                                   1  11         sum=12 
	//                                   11 11         sum=22

	//Case4: countA=3 -> 4 combinations: 1  1  1       sum=3 
	//                                   1  1  11      sum=13
	//                                   1  11 11      sum=23
	//                                   11 11 11      sum=33

	//Case5: countA=4 -> 5 combinations: 1  1  1  1    sum=4
	//                                   1  1  1  11   sum=14
	//                                   1  1  11 11   sum=24
	//                                   1  11 11 11   sum=34
	//                                   11 11 11 11   sum=44

	for _, cardTaken := range TakenCards{
		var val int
		switch(cardTaken[0]){
			case '0': val = 10
			case 'A': val = 1 //Firstly, we will consider all As as 1 and calculate the sum
			case 'J': val = 12
			case 'Q': val = 13
			case 'K': val = 14
			default: val,_ = strconv.Atoi(string(cardTaken[0]))
		}
		sum += val
	}

	//As the sum of each combination differs by 10, we will check countA+1 sums each greater by 10
	//If modul of any sum by HandCardValue is 0, player can take cards otherwise he can't
	for i:=0; i<=countA; i++{
		sum += i*10;
		if (sum%HandCardValue == 0){
			valid = true
			break
		}
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