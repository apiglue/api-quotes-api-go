package dataloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

const (
	redisMemberName = "quotes"
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

//Loadquotes - Loads data into redis
func Loadquotes() error {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		//return err
	}

	err = purgeData()
	if err != nil {
		fmt.Printf("Error purging redis: %s", err)
		log.Fatal(err)
		return err
	}

	err = loadData()
	if err != nil {
		//fmt.Printf("Error loading data into redis: %s", err)
		log.Fatal(err)
		return err
	}

	return nil
}

func purgeData() error {

	conn, err := redis.DialURL(os.Getenv("REDIS_URL"))
	if err != nil {
		// Handle error
		return err
	}

	conn.Do("DEL", redisMemberName)

	log.Print("Data purged")

	defer conn.Close()

	// conn, err := redis.Dial("tcp", os.Getenv("REDIS_SERVER"))
	// if err != nil {
	// 	return err
	// }

	// conn.Cmd("DEL quotes")

	// log.Print("Data purged")

	// defer conn.Close()

	return nil
}

func loadData() error {

	// Open our jsonFile
	jsonFile, err := os.Open("pkg/dataloader/db.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var quotes Quotes

	json.Unmarshal(byteValue, &quotes)

	// conn, err := redis.Dial("tcp", os.Getenv("REDIS_SERVER"))
	// if err != nil {
	// 	return err
	// }

	conn, err := redis.DialURL(os.Getenv("REDIS_URL"))
	if err != nil {
		// Handle error
		return err
	}

	for i := 0; i < len(quotes.Quotes); i++ {
		//conn.Cmd("SADD", redisMember, quotes.Quotes[i].Value)
		conn.Do("SADD", redisMemberName, quotes.Quotes[i].Value)
	}

	log.Print("Data loaded")

	defer conn.Close()

	return nil

}
