package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/initializers"
	"main.go/models"
)

var queue []models.User

func addPlayerHandler(c *gin.Context) {
	var newUser models.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(400, gin.H{"message": "Invalid user data"})
		return
	}
	result := initializers.DB.Create(&newUser)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}
	queue = append(queue, newUser)
	if len(queue) >= 2 {
		queue = queue[:2]
		c.JSON(201, gin.H{"message": "Waiting to start the game"})
		startGame(queue[0], queue[1], c)
	} else {
		c.JSON(201, gin.H{"message": "Wainting for the other player to join"})
	}
}
func startGame(player1 models.User, player2 models.User, c *gin.Context) {
	// alociram deck
	response, err := http.Get("https://www.deckofcardsapi.com/api/deck/new/shuffle/?deck_count=1")
	if err != nil {
		c.JSON(500, gin.H{"message": "Error starting the game"})
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var deckResponse struct {
			DeckID string `json:"deck_id"`
		}

		err = json.NewDecoder(response.Body).Decode(&deckResponse)
		if err != nil {
			c.JSON(500, gin.H{"message": "Error starting the game"})
			return
		}
		//pravim gejm
		newGame := models.Game{
			Score:         0,
			DeckPile:      deckResponse.DeckID,
			TablePile:     "table",
			HandPile:      "hand1",
			CollectedPile: "taken1",
			First:         true,
			UserID:        int(player1.ID),
			User:          player1,
		}

		result := initializers.DB.Create(&newGame)
		if result.Error != nil {
			c.JSON(500, gin.H{"message": "Error starting the game"})
			return
		}

		newGame2 := models.Game{
			Score:         0,
			DeckPile:      deckResponse.DeckID,
			TablePile:     "table",
			HandPile:      "hand2",
			CollectedPile: "taken2",
			First:         false,
			UserID:        int(player2.ID),
			User:          player2,
		}

		result2 := initializers.DB.Create(&newGame2)
		if result2.Error != nil {
			c.JSON(500, gin.H{"message": "Error starting the game"})
			return
		} else {
			c.JSON(201, gin.H{"message": "Game has started", "game1": newGame, "game2": newGame2})
		}

		// //uzimam karte iz spila
		// URL := fmt.Sprintf("https://www.deckofcardsapi.com/api/deck/%s/draw/?count=2", newGame.DeckPile)
		// drawedCards, err := http.Get(URL)
		// //dodjeljujem te karte u hand pile

	} else {
		c.JSON(500, gin.H{"message": "Error starting the game"})
	}
}

func InitializeHandlers(router *gin.Engine) {
	router.POST("/addPlayer", addPlayerHandler)
}
