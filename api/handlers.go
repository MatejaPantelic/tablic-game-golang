package api

import (
	"encoding/json"
	"fmt"
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

	fmt.Println(newDeckResponse.DeckId)
}

func InitializersHandlers(r *gin.Engine) {
	r.GET("/cards", newDeckHandler)
}
