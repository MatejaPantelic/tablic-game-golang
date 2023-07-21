package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/models"
)

func newDeckHandler(c *gin.Context) {

	resp, errURL := http.Get("https://www.deckofcardsapi.com/api/deck/new/")
	if errURL != nil {
		log.Fatal(errURL)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var newDeckResponse models.NewDeckResponse
	json.Unmarshal(body, &newDeckResponse)

	c.JSON(http.StatusOK, gin.H{
		"response": newDeckResponse})
}

func InitializersHandlers(r *gin.Engine) {
	r.GET("/cards", newDeckHandler)
	r.GET("/takecardsfromtable/", TakeCardsFromTable)
}
