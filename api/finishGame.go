package api

import (
	"github.com/gin-gonic/gin"
	"main.go/models"
	"log"
	"net/http"
	"fmt"
	"strings"
	"main.go/constants"
	"main.go/initializers"
)



func returnCardsToDeck(deckId string){
	url := fmt.Sprintf(constants.RETURN_TO_DECK_URL, deckId)
	_, errURL := http.Get(url)
	if errURL != nil {
		log.Fatal(errURL)
	}
}

func numberOfCards(deckId string) string {

	Player1TakenPile := Piles(deckId, "taken1")
	remainingPlayer1 := Player1TakenPile.Remaining

	Player2TakenPile := Piles(deckId, "taken2")
	remainingPlayer2 := Player2TakenPile.Remaining

	if remainingPlayer1 > remainingPlayer2 {
		return "taken1"
	} else if remainingPlayer1 == remainingPlayer2 {
		return "equal"
	} else {
		return "taken2"
	}

}

func FinishGame(c *gin.Context, deckId string){

	//Check if player's hands are empty
	//If one of hands is not empty round can continue
	if(notEmptyHands(deckId)){
		c.JSON(http.StatusOK, gin.H{"message": "You can continue round - next player's on turn"})
		return
	}

	//Check if deck is empty
	//If it isn't empty draw draw 2x6 cards from deck to player's hands
	if(!emptyDeck(deckId, "hand1")){
		createPile("6", deckId, "hand1", c)
		createPile("6", deckId, "hand2", c)
		return
	}

	//If deck is empty round is over
	//Check if table is empty
	if(!emptyTable(deckId)){
		//Draw cards from table
		cardsList := listPileCards(deckId, "table")
		cards := make([]string, 0)
		for _,card := range cardsList{
			cards = append(cards, card.Code)
		}
		cardsString := strings.Join(cards, ",")
		drawCardsFromPile(deckId, "table", cardsString)

		//Check who collected last
		var game models.Game
		result := initializers.DB.Model(&game).Where("deck_pile = ? AND collected_last = true", deckId).Find(&game)
		if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": result.Error})
		}

		//Add to taken pile of the player who collected last
		addToPile(deckId, game.CollectedPile, cardsString)
		
		//Update points
		Score(deckId, game.CollectedPile, cardsString, false)

		//Check who has more cards
		playerMoreCards := numberOfCards(deckId)

		//update points
		if(playerMoreCards != "equal"){
			var game models.Game
			err := initializers.DB.Where("deck_pile = ? and  collected_pile = ?", deckId, playerMoreCards).Find(&game).Error
			if err != nil {
				log.Fatal("Can't find game")
			}

			game.Score += 3

			result := initializers.DB.Model(&game).Where("collected_pile = ? AND deck_pile = ?", game.CollectedPile, deckId).Update("score", game.Score)
			if result.Error != nil {
				log.Fatal("Error updating score")
			}
		}
		
	}

	//Check points
	//If one of the players passed 100 finish game
	var game models.Game
	var games []models.Game
	result := initializers.DB.Model(&game).Where("deck_pile = ? AND collected_last = true", deckId).Find(&games)
	if result.Error != nil {
	c.JSON(http.StatusBadRequest, gin.H{"message": result.Error})
	}

	end := false
	for _,game := range games{
		if(game.Score > 100){
			end = true
			break
		}
	}

	if(end){
		// for _,game := range games{
			result :=initializers.DB.Model(&game).Where("deck_pile = ?", deckId).Update("game_finished", true)
			if result.Error != nil {
				c.JSON(http.StatusOK, gin.H{"message": result.Error})
			}
		// }
		return
	}

	//If game is not finshed - create new round
	//Move all cards from piles to deck
	returnCardsToDeck(deckId)
	createPile("6", deckId, "hand1", c)
	createPile("6", deckId, "hand2", c)
	createPile("4", deckId, "table", c)
}