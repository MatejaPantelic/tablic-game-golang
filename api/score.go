package api

import (
	"strings"
	"github.com/gin-gonic/gin"
	"main.go/initializers"
	"main.go/models"
)

func Score(deckId string, takenPile string, cards string, table bool, c *gin.Context){
	var game models.Game
	err := initializers.DB.Where("deck_pile = ? and  collected_pile = ?", deckId, takenPile).Find(&game).Error
	errorCheck(err, 400,"Error during connecting to base",c)

	var oldScore = game.Score
	var newScore = 0

	//calculating points for every taken card
	codes := strings.Split(cards, ",")
	for _, code := range codes {
		newScore += calculateCardPoints(code)
	}

	if(table && emptyTable(deckId)){
		newScore++
	}

	game.Score = oldScore + newScore

	result := initializers.DB.Model(&game).Where("collected_pile = ? AND deck_pile = ?", takenPile, deckId).Update("score", game.Score)
	errorCheck(result.Error, 400,"Cannot update score",c)
	
}