package main

import (
	"os"

	"github.com/acheong08/ShareGPT/checks"
	"github.com/acheong08/ShareGPT/typings"
	gin "github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var (
	redisAddr = os.Getenv("REDIS_ADDRESS")
	redisPass = os.Getenv("REDIS_PASSWORD")
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
	})
}

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
		if creditSummary.Error.Message != "" {
			c.JSON(400, gin.H{
				"error": creditSummary.Error.Message,
			})
			return
		}
		if creditSummary.TotalAvailable < 0.1 {
			c.JSON(400, gin.H{
				"error": "Not enough credits",
			})
			return
		}
		// Return credit summary
		c.JSON(200, creditSummary)
		// Save to Redis
		go func(creditSummary typings.CreditSummary) {
			// Save to Redis
			rdb.Set(api_key.APIKey, creditSummary.TotalAvailable, 0)
		}(creditSummary)
	})
	router.Run()
}
