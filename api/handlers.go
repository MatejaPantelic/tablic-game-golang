package api

import (
	"encoding/json"
	//"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"main.go/initializers"
	"main.go/models"
)

var queue []models.User

func addPlayerHandler(c *gin.Context) {
	var newUser models.User

	err := c.BindJSON(&newUser)
	if err != nil {
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

		//Uzimanje karata iz spila, formiranje karata za prvog igraca
		URL := "https://www.deckofcardsapi.com/api/deck/"+deckResponse.DeckID+"/draw/?count="
		URL1 := URL+"6"
		drawnCards, err := http.Get(URL1)
		if err != nil {
			c.JSON(500, gin.H{"message": "Error starting the game"})
			return
		}
		defer drawnCards.Body.Close()
		
		//Karte se dodjeljuju u hand pile
		if drawnCards.StatusCode == http.StatusOK{

			var drawnCardsResponse struct {
				Success   bool   `json:"success"`
				DeckId    string `json:"deck_id"`
				Cards     []models.Card `json:"cards"`
				Remaining int    `json:"remaining"`
			}
			
			err = json.NewDecoder(drawnCards.Body).Decode(&drawnCardsResponse)
			if err != nil {
				c.JSON(500, gin.H{"message": "Error starting the game"})
				return
			}

			//Making of first player hand
			cardCodes := ""
            for i := 0; i < len(drawnCardsResponse.Cards); i++ {
				cardCodes+=drawnCardsResponse.Cards[i].Code
				if(i < len(drawnCardsResponse.Cards)-1){
					cardCodes+=","
				}
			}

			URL_AddToPiles := "https://deckofcardsapi.com/api/deck/"+deckResponse.DeckID+"/pile/"+newGame.HandPile+"/add/?cards="+cardCodes
			player1Hand, err := http.Get(URL_AddToPiles)
			if err != nil {
		        c.JSON(500, gin.H{"message": "Error starting the game"})
		        return
	        }
	        defer player1Hand.Body.Close()

			if player1Hand.StatusCode == http.StatusOK{
				var player1HandPileResponse models.AddingToPilesResponse
				
				err = json.NewDecoder(player1Hand.Body).Decode(&player1HandPileResponse)
				if err != nil {
					c.JSON(500, gin.H{"message": "Error starting the game"})
					return
				}
			}
		}
		/*-------------------------------------------------------------------------------------------*/
				//Uzimanje karata iz spila, formiranje karata za drugog igraca
				URL2 := URL+"6"
				drawnCards2, err := http.Get(URL2)
				if err != nil {
					c.JSON(500, gin.H{"message": "Error starting the game"})
					return
				}
				defer drawnCards2.Body.Close()
				
				//Karte se dodjeljuju u hand pile
				if drawnCards2.StatusCode == http.StatusOK{
		
					var drawnCards2Response struct {
						Success   bool   `json:"success"`
						DeckId    string `json:"deck_id"`
						Cards     []models.Card `json:"cards"`
						Remaining int    `json:"remaining"`
					}
					
					err = json.NewDecoder(drawnCards2.Body).Decode(&drawnCards2Response)
					if err != nil {
						c.JSON(500, gin.H{"message": "Error starting the game"})
						return
					}
		
					//Making of second player hand
					cardCodes := ""
					for i := 0; i < len(drawnCards2Response.Cards); i++ {
						cardCodes+=drawnCards2Response.Cards[i].Code
						if(i < len(drawnCards2Response.Cards)-1){
							cardCodes+=","
						}
					}
		
					URL_AddToPiles := "https://deckofcardsapi.com/api/deck/"+deckResponse.DeckID+"/pile/"+newGame2.HandPile+"/add/?cards="+cardCodes
					player2Hand, err := http.Get(URL_AddToPiles)
					if err != nil {
						c.JSON(500, gin.H{"message": "Error starting the game"})
						return
					}
					defer player2Hand.Body.Close()
		
					if player2Hand.StatusCode == http.StatusOK{
						var player2HandPileResponse models.AddingToPilesResponse
		
						err = json.NewDecoder(player2Hand.Body).Decode(&player2HandPileResponse)
						if err != nil {
							c.JSON(500, gin.H{"message": "Error starting the game"})
							return
						}
					}
				}
				/*-----------------------------------------------------------------------------------------------------------*/
				//Uzimanje karata iz spila, formiranje karata za talon
				URL3 := URL+"4"
				drawnCardsTable, err := http.Get(URL3)
				if err != nil {
					c.JSON(500, gin.H{"message": "Error starting the game"})
					return
				}
				defer drawnCardsTable.Body.Close()
				
				//Karte se dodjeljuju u hand pile
				if drawnCardsTable.StatusCode == http.StatusOK{
		
					var drawnCardsTableResponse struct {
						Success   bool   `json:"success"`
						DeckId    string `json:"deck_id"`
						Cards     []models.Card `json:"cards"`
						Remaining int    `json:"remaining"`
					}
					
					err = json.NewDecoder(drawnCardsTable.Body).Decode(&drawnCardsTableResponse)
					if err != nil {
						c.JSON(500, gin.H{"message": "Error starting the game"})
						return
					}
		
					//Making of table cards
					cardCodes := ""
					for i := 0; i < len(drawnCardsTableResponse.Cards); i++ {
						cardCodes+=drawnCardsTableResponse.Cards[i].Code
						if(i < len(drawnCardsTableResponse.Cards)-1){
							cardCodes+=","
						}
					}
		
					URL_AddToPiles := "https://deckofcardsapi.com/api/deck/"+deckResponse.DeckID+"/pile/"+newGame.TablePile+"/add/?cards="+cardCodes
					Table, err := http.Get(URL_AddToPiles)
					if err != nil {
						c.JSON(500, gin.H{"message": "Error starting the game"})
						return
					}
					defer Table.Body.Close()
		
					if Table.StatusCode == http.StatusOK{
						var TablePileResponse models.AddingToPilesResponse
		
						err = json.NewDecoder(Table.Body).Decode(&TablePileResponse)
						if err != nil {
							c.JSON(500, gin.H{"message": "Error starting the game"})
							return
						}
					}
				}		

	} else {
		c.JSON(500, gin.H{"message": "Error starting the game"})
	}
}

func InitializeHandlers(router *gin.Engine) {
	router.POST("/addPlayer", addPlayerHandler)
}
