package main

import (
	"fmt"
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
			err := rdb.Set(api_key.APIKey, creditSummary.TotalAvailable, 0).Err()
			if err != nil {
				println(fmt.Errorf("error saving to Redis: %v", err))
			}
		}(creditSummary)
	})
	router.GET("/api_key/random", func(c *gin.Context) {
		// Check authentication
		if c.GetHeader("Authorization") != os.Getenv("AUTHORIZATION") {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			return
		}
		// Get random API key from Redis
		key, err := rdb.RandomKey().Result()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		if key == "" {
			c.JSON(400, gin.H{
				"error": "No API keys",
			})
			return
		}
		// Get credit summary
		creditSummary, err := checks.GetCredits(key)
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
		c.JSON(200, gin.H{
			"api_key":        key,
			"credit_summary": creditSummary,
		})
	})
	router.POST("/api_key/delete", func(c *gin.Context) {
		// Delete API key from Redis
		var api_key typings.APIKeySubmission
		err := c.BindJSON(&api_key)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Error binding JSON",
			})
			return
		}
		if api_key.APIKey == "" {
			c.JSON(400, gin.H{
				"error": "API key is empty",
			})
			return
		}
		err = rdb.Del(api_key.APIKey).Err()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "API key deleted",
		})
	})
	router.Run()
}
