package api

import (
	"github.com/gin-gonic/gin"
	"main.go/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"fmt"
	// "bufio"
	// "os"
	// "strings"
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

	// var Piles []models.Players
	// json.Unmarshal(body, &ListCardResponse.Piles.Players)

	//Kako da iscitam sve karte?
	for _, player := range ListCardResponse.Piles.Players {
		fmt.Println("Remaining:", player.Remaining)

		// for _, card := range player.Cards {
		// 	fmt.Println("Card:", card)
		// }
	}

	fmt.Println(len(ListCardResponse.Piles.Players))

	return

	// c.JSON(http.StatusOK, gin.H{
	// 	"response": ListCardsResponse})
}


func TakeCardsFromTable(c *gin.Context){
	//choose card from hand
	// scanner := bufio.NewScanner(os.Stdin)
	// fmt.Print("Which card you want to throw:")
	// scanner.Scan()
	// HandCard := scanner.Text()
	// fmt.Println(HandCard)

	//validate card from hand
	// var HandCards []models.CardList = listPileCards("y54h3qiktc1l", "hand1")
	listPileCards("y54h3qiktc1l", "hand1")

	//choose cards from table
	// var stringsList []string
	// fmt.Println("Enter cards you want to take from table (type 'done' to finish):")
	// for {
	// 	scanner.Scan()
	// 	card := scanner.Text()
	// 	if strings.ToLower(card) == "done" {
	// 		break
	// 	}
	// 	CardsList = append(CardsList, strings.ToUpper(card))
	// }
	// fmt.Println("Cards you want to take:")
	// for _, card := range CardsList {
	// 	fmt.Println(cards.Code)
	// }

	//validate cards from table
	// var TableCards []models.CardList = listPileCards("y54h3qiktc1l", "hand1")
	listPileCards("y54h3qiktc1l", "table")

	//check values!

	//if valid
	//take card from hand

	//take card(s) from table

	//move cards to taken

}