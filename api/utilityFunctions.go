package api

import(
	"net/http"
	"io"
	"log"
	"main.go/models"
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

//Function for checking if cards exists in specific pile
func existsInPile(cardCode string, pile []models.CardList)(exist bool){
	exist = false
	for _, card := range pile {
		if card.Code == cardCode{
			exist=true
		}
	}
	return
}