package main

import (
	"log"
	"net/http"
	"os"

	"github.com/apiglue/api-quotes-api-go/pkg/dataloader"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mediocregopher/radix.v2/redis"
)

const (
	redisMember = "quotes"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	err = dataloader.Loadquotes()
	if err != nil {
		//log.Fatal("Dataloader thrown an error: %s", err)
		return
	}

	//gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	v1 := router.Group("/")
	{
		v1.GET("/random", getRandomQuote)
	}
	router.Run()

}

//GetRandomQuote - GET A RANDOM QUOTE
func getRandomQuote(c *gin.Context) {

	conn, err := redis.Dial("tcp", os.Getenv("REDIS_SERVER"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	conn.Cmd("SADD", redisMember)

	quote, err := conn.Cmd("SRANDMEMBER", redisMember).Str()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, quote)

	defer conn.Close()

}
