package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Request struct {
	CadastralNumber string  `json:"cadastral_number"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
}

type Response struct {
	Result bool    `json:"result"`
	Delay  float64 `json:"delay"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	delayMaxStr := os.Getenv("DELAY_MAX")
	delayMax := 60
	if delayMaxStr != "" {
		if val, err := strconv.Atoi(delayMaxStr); err == nil {
			delayMax = val
		}
	}

	router := gin.Default()

	router.POST("/api/result", func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Received request for cadastral number: %s", req.CadastralNumber)

		// imitate a proccess (from 1 to delayMax seconds)
		delay := rand.Intn(delayMax) + 1
		time.Sleep(time.Duration(delay) * time.Second)

		// random result
		result := rand.Intn(2) == 1

		response := Response{
			Result: result,
			Delay:  float64(delay),
		}

		log.Printf("Processed request: result=%v, delay=%ds", result, delay)

		c.JSON(http.StatusOK, response)
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "mock server is running"})
	})

	log.Printf("Mock server starting on port %s", port)
	log.Fatal(router.Run(":" + port))
}
