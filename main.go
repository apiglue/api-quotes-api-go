package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mediocregopher/radix.v2/redis"
)

const (
	redisMember = "quotes"
)

// Quotes struct which contains
// an array of quotes
type Quotes struct {
	Quotes []Quote `json:"quotes"`
}

//Quote struct
// single quote
type Quote struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

func main() {

	// Open our jsonFile
	jsonFile, err := os.Open("db.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Successfully Opened db.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var quotes Quotes

	json.Unmarshal(byteValue, &quotes)

	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(quotes.Quotes); i++ {
		conn.Cmd("SADD", redisMember, quotes.Quotes[i].Value)
	}

	defer conn.Close()

}

// // LOAD data into redis
// func populateQuotes() error {
// 	var err error

// 	return nil
// }
