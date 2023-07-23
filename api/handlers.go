package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/constants"
	"main.go/models"
)

func newDeckHandler(c *gin.Context) {
	resp, err := http.Get(constants.NEW_DECK_URL)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"messege": err,
		})
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"messege": err,
		})
	}

	var newDeckResponse models.NewDeckResponse
	json.Unmarshal(body, &newDeckResponse)
	c.JSON(http.StatusOK, gin.H{
		"response": newDeckResponse})
}

func InitializersHandlers(r *gin.Engine) {
	r.GET("/cards", newDeckHandler)
	r.GET("/takecardsfromtable/:deckId/:handPile/:takenPile", TakeCardsFromTable)
	r.GET("/throwCard/:cardCode/:deckId/:playerPile", throwCardHandler)
}
