package api

import (
	"github.com/gin-gonic/gin"
	// "tablic/models"
	"net/http"
	"fmt"
	"log"
)

func drowFromPileHandler(){
	// res, err := http.Get("https://www.deckofcardsapi.com/api/deck/<<deck_id>>/pile/<<pile_name>>/draw/?cards=AS")
}

func addToPileHandler(){
	// res, err := http.Get("https://www.deckofcardsapi.com/api/deck/<<deck_id>>/pile/<<pile_name>>/add/?cards=AS,2S")
}

func TakeCardsFromTable(c *gin.Context){
	//take card from hand
	urlFromHand := fmt.Sprintf("/drawfrompile/%s/%s/%s", deck, pile, cards)
	resFromHand, errFromHand := http.Get(urlFromHand)

	//validate card
	if(errFromHand != nil){
		log.Fatal(errFromHand)
	}

	//take card(s) from table
	urlFromTable := fmt.Sprintf("/drawfrompile/%s/%s/%s", deck, pile, cards)
	resFromTable, errFromTable := http.Get(urlFromTable)

	//validate cards
	if(errFromTable != nil){
		log.Fatal(errFromTable)
	}

	//check values!

	//if values are ok, move cards to taken
	urlToTaken := fmt.Sprintf("/addtopile/%s/%s/%s", deck, pile, cards)
	resToTaken, errToTaken:= http.Get(urlToTaken)

	// //if not, move them back to hand and table
	// urlToHand := fmt.Sprintf("/drawfrompile/%s/%s/%s", deck, pile, cards)
	// resToHand, errToHand := http.Get(urlToHand)

	// urlToTable := fmt.Sprintf("/drawfrompile/%s/%s/%s", deck, pile, cards)
	// resToTable, errToTable := http.Get(urlToTable)

}

func InitializeHandlers(router *gin.Engine){
	router.GET("/drawfrompile/:deck/:pile/:cards", drowFromPileHandler)
	router.GET("/addtopile/:deck/:pile/:cards", addToPileHandler)

	router.GET("/takecardsfromtable", TakeCardsFromTable)
}