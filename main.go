package main

import (
	"fmt"
	"io"
	"os"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"

	"github.com/acheong08/ShareGPT/checks"
	"github.com/acheong08/ShareGPT/typings"
	gin "github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var (
	jar     = tls_client.NewCookieJar()
	options = []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(360),
		tls_client.WithClientProfile(tls_client.Chrome_110),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}
	client, _ = tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
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
	router.GET("/ping", func(c *gin.Context) {
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
			// Save to Redis without expiration
			err := rdb.Set(api_key.APIKey, creditSummary.TotalAvailable, 0).Err()
			if err != nil {
				println(fmt.Errorf("error saving to Redis: %v", err))
			}
		}(creditSummary)
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
	router.POST("/v1/chat", proxy)
	router.Run()
}

func proxy(c *gin.Context) {

	var url string
	var err error
	var request_method string
	var request *http.Request
	var response *http.Response

	url = "https://api.openai.com/v1/chat/completions"
	request_method = c.Request.Method

	request, err = http.NewRequest(request_method, url, c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	request.Header.Set("Host", "api.openai.com")
	request.Header.Set("Origin", "https://platform.openai.com/playground")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Keep-Alive", "timeout=360")
	request.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	// Authorization
	var authorization string
	if c.Request.Header.Get("Authorization") == "" {
		// Set authorization header from Redis
		random_key, err := rdb.RandomKey().Result()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get random key from Redis"})
			println(err.Error())
			return
		}
		authorization = "Bearer " + random_key
	} else {
		// Set authorization header from request
		authorization = c.Request.Header.Get("Authorization")
	}
	request.Header.Set("Authorization", authorization)

	response, err = client.Do(request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer response.Body.Close()
	c.Header("Content-Type", response.Header.Get("Content-Type"))
	// Get status code
	c.Status(response.StatusCode)
	c.Stream(func(w io.Writer) bool {
		// Write data to client
		io.Copy(w, response.Body)
		return false
	})

}
