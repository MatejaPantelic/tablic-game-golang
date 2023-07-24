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
	"main.go/initializers"
)

//Function for listing cards in a pile
func listPileCards(deck string, pileName string)(CardsArray []models.CardList){
	url := fmt.Sprintf(constants.LIST_PILE_CARDS_URL, deck, pileName)
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
	url := fmt.Sprintf(constants.DRAW_CARDS_FROM_PILE_URL, deck, pileName, cards)
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
	url := fmt.Sprintf(constants.ADD_TO_PILE_URL, deck, pileName, cards)
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

//Function for checking if sum of cards group is same as hand card value
func isGroupValid(c *gin.Context, HandCard string, TakenCards []string)(valid bool){
	valid = false

	//CHECK VALUES
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
	fmt.Println(HandCardValue)

	//As the sum of each combination differs by 10, we will check countA+1 sums each greater by 10
	//If modul of any sum by HandCardValue is 0, player can take cards otherwise he can't
	for i:=0; i<=countA; i++{
		sum += i*10;
		if (sum == HandCardValue){
			valid = true
			return
		}
	}

	return
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
		c.JSON(http.StatusBadRequest, gin.H{"response": "Failed to read body"})
		return
	}
	HandCard := RequestData.HandCard
	TakenCardsString := RequestData.TakenCards
	TakenCardsGroups := strings.Split(TakenCardsString, ";")
	var TakenCards []string

	var valid bool = true

	//VALIDATE EACH CARDS GROUP
	for _, group := range TakenCardsGroups{
		TakenCards = strings.Split(group, ",")

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

		if(!isGroupValid(c, HandCard, TakenCards)){
			valid = false
			break
		}

	}

	//IF ONE OF THE CARDS GROUP IS NOT VALID
	if(!valid){
		c.JSON(http.StatusNotFound, gin.H{"response": "You can't take chosen cards"})
		return
	}

	//IF VALID MOVE CARDS FROM HAND AND TABLE PILE TO TAKEN PILE
	drawCardsFromPile(deckId, handPile, HandCard)
	separator := ","
	cards := strings.Join(TakenCards, separator)
	drawCardsFromPile(deckId, "table", cards)
	addToPile(deckId, takenPile, cards+","+HandCard)

	//NOTE THAT THIS PLAYER HAS COLLECTED CARDS LAST
	var game models.Game
	if(handPile == "hand1"){
		initializers.DB.Model(&game).Where("deck_pile = ? AND hand_pile = ?", deckId, "hand1").Update("collected_last", true) 
		initializers.DB.Model(&game).Where("deck_pile = ? AND hand_pile = ?", deckId, "hand2").Update("collected_last", false)
	}else{
		initializers.DB.Model(&game).Where("deck_pile = ? AND hand_pile = ?", deckId, "hand2").Update("collected_last", true)
		initializers.DB.Model(&game).Where("deck_pile = ? AND hand_pile = ?", deckId, "hand1").Update("collected_last", false)
	}

	c.JSON(http.StatusOK, gin.H{"response": "Cards are moved from hand and table pile to taken pile"})
	
}