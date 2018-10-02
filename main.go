package main

import (
	"log"
	"net/http"
	"os"

	"github.com/apiglue/api-quotes-api-go/pkg/dataloader"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

const (
	redisMember = "quotes"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	err := dataloader.Loadquotes()
	if err != nil {
		//log.Fatal("Dataloader thrown an error: %s", err)
		return
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	v1 := router.Group("/")
	{
		v1.GET("/random", getRandomQuote)
	}
	router.Run()

	log.Print("Server started on port: " + port)

}

//GetRandomQuote - GET A RANDOM QUOTE
func getRandomQuote(c *gin.Context) {

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("$PORT must be set")
	}

	conn, err := redis.DialURL(redisURL)
	if err != nil {
		// Handle error
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	quote, err := redis.String(conn.Do("SRANDMEMBER", redisMember))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, quote)

	defer conn.Close()

}
