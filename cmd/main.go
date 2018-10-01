package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/apiglue/api-quotes-api-go/pkg/dataloader"
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()
	router.HandleFunc("/random", GetRandomQuote).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+port, router))
}

//GetRandomQuote - GET A RANDOM QUOTE
func GetRandomQuote(w http.ResponseWriter, r *http.Request) {

	conn, err := redis.Dial("tcp", os.Getenv("REDIS_SERVER"))
	if err != nil {
		return
	}

	conn.Cmd("SADD", redisMember)

	quote, err := conn.Cmd("SRANDMEMBER", redisMember).Str()
	if err != nil {
		// handle err
	}

	json.NewEncoder(w).Encode(quote)

	defer conn.Close()

}
