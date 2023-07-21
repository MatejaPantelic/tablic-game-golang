package api

import (
	"github.com/gin-gonic/gin"
	"main.go/models"
	// "main.go/initializers"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
)

func listPileCards(deck string, pileName string)(Cards []models.CardList){
	url := fmt.Sprintf("https://www.deckofcardsapi.com/api/deck/%s/pile/%s/list/", deck, pileName)
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

	var CardsArray []models.CardList
	switch(pileName){
	case "hand1": CardsArray = ListCardResponse.Piles.Hand1.Cards
	case "hand2": CardsArray = ListCardResponse.Piles.Hand2.Cards
	case "table": CardsArray = ListCardResponse.Piles.Table.Cards
	default: 
	}

	for _, card := range CardsArray {
		Cards = append(Cards, card)
	}

	return
}

func drawCardsFromPile(deck string, pileName string, cards string){
	url := fmt.Sprintf("https://www.deckofcardsapi.com/api/deck/%s/pile/%s/draw/?cards=%s", deck, pileName, cards)
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

func addToPile(deck string, pileName string, cards string){
	url := fmt.Sprintf("https://www.deckofcardsapi.com/api/deck/%s/pile/%s/add/?cards=%s", deck, pileName, cards)
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


func TakeCardsFromTable(c *gin.Context){
	//EXTRACT PARAMETERS
	// playerId := c.Param("player")
	// deckId := c.Param("deck")
	// var game models.Game
	// initializers.DB.Where("user_id=? AND deck_pile=?", playerId, deckId).Find(&game)

	//CHOOSE CARD FROM HAND
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Which card you want to throw:")
	scanner.Scan()
	HandCard := scanner.Text()
	fmt.Println(HandCard)

	//VALIDATE CARD FROM HAND
	// var HandCards []models.CardList = listPileCards(game.DeckId, game.HandPile)
	var HandCards []models.CardList = listPileCards("y54h3qiktc1l", "hand1")
	var existInHand bool = false
	for _, card := range HandCards {
		if card.Code == HandCard{
			existInHand=true
		}
	}
	if(!existInHand){
		c.JSON(http.StatusOK, gin.H{"response": "The selected card is not in your hand."})
	}

	//CHOOSE CARDS FROM TABLE
	var CardsList []string
	fmt.Println("Enter cards you want to take from table (type 'done' to finish):")
	for {
		scanner.Scan()
		card := scanner.Text()
		if strings.ToLower(card) == "done" {
			break
		}
		CardsList = append(CardsList, strings.ToUpper(card))
	}
	// fmt.Println("Cards you want to take:")
	// for _, card := range CardsList {
	// 	fmt.Println(card)
	// }

	//VALIDATE CARDS FROM TABLE
	// var TableCards []models.CardList = listPileCards(game.DeckId, game.TablePile)
	var TableCards []models.CardList = listPileCards("y54h3qiktc1l", "table")
	for _, cardT := range CardsList {
		existInHand = false
		for _, cardTable := range TableCards{
			if cardT == cardTable.Code{
				existInHand = true
			}
		}
		if(existInHand == false)	{
			c.JSON(http.StatusOK, gin.H{"response": "Some of selected cards is not on the table."})
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
	}
	for _, cardT := range CardsList{
		valStr := cardT[0]
		var val int
		switch(cardT[0]){
			case 'J': val = 12
			case 'Q': val = 13
			case 'K': val = 14
			default: val,_ = strconv.Atoi(string(valStr))
		}
		sum += val
		fmt.Println(val, sum)
	}
	if (sum%HandCardValue == 0){
		valid = true
	}
	if(valid){
		// drawCardsFromPile(game.DeckId, game.HandPile, HandCard)
		drawCardsFromPile("y54h3qiktc1l", "hand1", HandCard)
		separator := ","
		cards := strings.Join(CardsList, separator)
		// drawCardsFromPile(game.DeckId, game.TablePile, cards)
		drawCardsFromPile("y54h3qiktc1l", "hand1", cards)
		// addToPile((game.DeckId, game.TakePile, cards+","+HandCard)
		addToPile("y54h3qiktc1l", "taken1", cards+","+HandCard)

		c.JSON(http.StatusOK, gin.H{"response": "Cards are moved from hand and table to taken pile"})
	}else{
		c.JSON(http.StatusOK, gin.H{"response": "You can't take chosen cards"})
	}

}