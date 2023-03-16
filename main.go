package main

import (
	"github.com/acheong08/ShareGPT/checks"
	"github.com/acheong08/ShareGPT/typings"
	gin "github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Any("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/api_key/submit", func(c *gin.Context) {
		var api_key typings.APIKeySubmission
		c.BindJSON(&api_key)
		if api_key.APIKey == "" {
			c.JSON(400, gin.H{
				"error": "API key is empty",
			})
			return
		}
		// Check if API key is valid
		creditSummary, err := checks.GetCredits(api_key.APIKey)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Return credit summary
		c.JSON(200, creditSummary)
	})
	router.Run()
}
