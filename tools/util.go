package tools

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/models"
)

func JsonParse(reponse *http.Response, c *gin.Context) (*models.DrawResponse, error) {
	//part of code for accepting response
	defer reponse.Body.Close()
	body, _ := io.ReadAll(reponse.Body)

	//variable for store reponse from http request in acceptable format
	var drawResponse models.DrawResponse

	errHandUm := json.Unmarshal(body, &drawResponse)

	// Error checking for JSON unmarshaling.
	if errHandUm != nil {
		return nil, errHandUm

	}
	return &drawResponse, nil
}
