package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"

	"github.com/acheong08/ShareGPT/checks"
	"github.com/acheong08/ShareGPT/typings"
	"github.com/fvbock/endless"
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
	api_keys  []string
)

var rdb *redis.Client

func init() {
	if redisAddr == "" {
		panic("REDIS_ADDRESS is not set")
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
	})
	// Fetch all keys from redis
	keys, err := rdb.Keys("*").Result()
	if err != nil {
		panic(err)
	}
	api_keys = keys
}

func cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func assign_api_keys() {
	// Fetch all keys from redis
	keys, err := rdb.Keys("*").Result()
	if err != nil {
		return
	}
	api_keys = keys
}

func main() {
	router := gin.Default()
	router.Use(cors)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// OPTIONS any route
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
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
		if creditSummary.HardLimitUSD < 1 {
			c.JSON(400, gin.H{
				"error": "Not enough credits",
			})
			return
		}
		// Return credit summary
		c.JSON(200, creditSummary)
		// Save to Redis
		go func(creditSummary typings.BillingSubscription) {
			// Save to Redis without expiration
			err := rdb.Set(api_key.APIKey, creditSummary.HardLimitUSD, 0).Err()
			if err != nil {
				println(fmt.Errorf("error saving to Redis: %v", err))
			}
		}(creditSummary)

		assign_api_keys()
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
		assign_api_keys()

	})
	// Status check (takes random API key from Redis and returns its credit summary)
	router.GET("/api_key/status", func(c *gin.Context) {
		// Get random API key from Redis
		key, err := rdb.RandomKey().Result()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
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
		c.JSON(200, creditSummary)
		assign_api_keys()
	})

	router.POST("/v1/chat", proxy)
	HOST := os.Getenv("HOST")
	if HOST == "" {
		HOST = "127.0.0.1"
	}
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8082"
	}
	endless.ListenAndServe(HOST+":"+PORT, router)
}

func proxy(c *gin.Context) {

	var url string
	var err error
	var request_method string
	var request *http.Request
	var response *http.Response

	url = "https://api.openai.com/v1/chat/completions"
	request_method = c.Request.Method

	var body []byte
	body, err = io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Convert to JSON
	var jsonBody map[string]interface{}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Set model to `gpt-3.5-turbo`
	jsonBody["model"] = "gpt-3.5-turbo"
	// Convert back to bytes
	body, err = json.Marshal(jsonBody)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	request, err = http.NewRequest(request_method, url, bytes.NewReader(body))
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
	var random_key string
	if c.Request.Header.Get("Authorization") == "" {
		// Choose random API key from api_keys array
		random_key = api_keys[rand.Intn(len(api_keys))]
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
	if response.StatusCode == 401 {
		// Delete API key from Redis
		err = rdb.Del(random_key).Err()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Remove API key from api_keys array
		assign_api_keys()
	}
	if response.StatusCode == 429 {
		// Check HardLimitUsd
		creditSummary, err := checks.GetCredits(random_key)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		if creditSummary.HardLimitUSD == 0 {
			// Remove from Redis
			err = rdb.Del(random_key).Err()
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}
			assign_api_keys()
		}
	}
	c.Header("Content-Type", response.Header.Get("Content-Type"))
	// Get status code
	c.Status(response.StatusCode)
	c.Stream(func(w io.Writer) bool {
		// Write data to client
		io.Copy(w, response.Body)
		return false
	})

}
