package tools

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JsonParse(response *http.Response, c *gin.Context, parameter interface{}) error {
	//part of code for accepting response
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	//check error during reading response
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error during reading response!"})
	}

	// error expectation if exist
	errHandUm := json.Unmarshal(body, parameter)

	return errHandUm
}

func CheckError(code int, c *gin.Context, err error, message string) {

	if err != nil {
		if message != "" {
			c.JSON(code, gin.H{"message": message})
		} else {
			c.JSON(code, gin.H{"message": err})
		}
	}
}
