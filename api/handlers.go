package api

import (
	"encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"
	"main.go/constants"
	"main.go/initializers"
	"main.go/models"
	"fmt"
)

// Making the queue where player will be stored until game starts
var queue []models.User

//Helper function for error responses
func errorCheck( err error, errCode int,errMsg string ,c *gin.Context ){
	if err != nil {
		c.JSON(errCode, gin.H{"message": errMsg})
		return
	}
}

// Function- adds player to queue.
// Creates a new user according to model, checks for possible errors, adds user to database and to queue
// If enough players are in queue, remove those player from queue and calls function startGame
func addPlayerHandler(c *gin.Context) {
	var newUser models.User

	err := c.BindJSON(&newUser)
	errorCheck(err,400,"Invalid user data",c)

	result := initializers.DB.Create(&newUser)
	errorCheck(result.Error,http.StatusBadRequest,"Failed to create user",c)

	queue = append(queue, newUser)
	if len(queue) >= 2 {
		startGame(queue[0], queue[1], c)
		queue = queue[:2]
	} else {
		c.JSON(201, gin.H{"message": "Wainting for the other player to join"})
	}
}

// Called by addPlayerHandler
// If enough players are present, starts the game
// Creates a game, deck and piles (player hands and table cards)
func startGame(player1 models.User, player2 models.User, c *gin.Context) {
	// Deck creation/alocation
	errorResponse := "Error starting the game"
	response, err := http.Get(constants.NEW_SHUFFLED_DECK)
	errorCheck(err,http.StatusBadRequest,errorResponse,c)
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var deckResponse struct {
			DeckID string `json:"deck_id"`
		}

		err = json.NewDecoder(response.Body).Decode(&deckResponse)
		errorCheck(err,http.StatusBadGateway,errorResponse,c)

		//Game creation
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
		errorCheck(result.Error,500,errorResponse,c)
		
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
			c.JSON(500, gin.H{"message": errorResponse})
			return
		} else {
			c.JSON(201, gin.H{"message": "Game has started", "game1": newGame, "game2": newGame2})
		}

		//Taking 6 cards from deck and forming cards for player1 hands
		createPile("6",deckResponse.DeckID,newGame.HandPile,c)
		
		//Taking 6 cards from deck and forming cards for player2 hands
		createPile("6",deckResponse.DeckID,newGame2.HandPile,c)

		//Taking 4 cards from deck and forming cards for table
		createPile("4",deckResponse.DeckID,newGame.TablePile,c)

	} else {
		c.JSON(500, gin.H{"message": errorResponse})
	}
}

//Used to create 3 piles at the start of game
//Those 3 piles are player hands(for each player) and table pile
func createPile(numberOfCards string, deckID string, pileToBeCreated string,c *gin.Context) {
	draw_A_Card := fmt.Sprintf(constants.DRAW_A_CARD_URL, deckID, numberOfCards)
	drawnCards, err := http.Get(draw_A_Card)
	errorMessage:= "Error starting the game"
	errorCheck(err,500,errorMessage,c)
	defer drawnCards.Body.Close()
			
			
	if drawnCards.StatusCode == http.StatusOK{
	
		var drawnCardsResponse struct {
			Success   bool   `json:"success"`
			DeckId    string `json:"deck_id"`
			Cards     []models.Card `json:"cards"`
			Remaining int    `json:"remaining"`
		}
				
		err = json.NewDecoder(drawnCards.Body).Decode(&drawnCardsResponse)
		errorCheck(err,500,errorMessage,c)
	
		//Taking card codes for API URL
		cardCodes := ""
		for i := 0; i < len(drawnCardsResponse.Cards); i++ {
			cardCodes+=drawnCardsResponse.Cards[i].Code
			if(i < len(drawnCardsResponse.Cards)-1){
				cardCodes+=","
			}
		}
	
		addToPileURL := fmt.Sprintf(constants.ADD_TO_PILE_URL, deckID,pileToBeCreated,cardCodes)
		newPile, err := http.Get(addToPileURL)
		errorCheck(err,500,errorMessage,c)
		defer newPile.Body.Close()
	
		if newPile.StatusCode == http.StatusOK{
			var player1HandPileResponse models.AddingToPilesResponse
					
			err = json.NewDecoder(newPile.Body).Decode(&player1HandPileResponse)
			if err != nil {
				c.JSON(500, gin.H{"message": "Error starting the game"})
				return
			}
		}
	}	
}

func InitializeHandlers(router *gin.Engine) {
	router.POST("/addPlayer", addPlayerHandler)
}