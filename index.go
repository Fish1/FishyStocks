package main

import (
	"os"
	"io"
	"io/ioutil"
	"fmt"
	"time"
	"encoding/json"
	"net/http"
	"math/rand"
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
)

type Key struct {
	Public string `json:"public"`
	Secret string `json:"secret"`
}

type WatchList struct {
	Stocks [] string `json:"stocks"`
}

var watchList WatchList 

func init() {
	fmt.Println("Welcome to Fishy Stocks") 

	/* Read keys from file */
	var key Key
	file, err := os.Open("key.json")
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(bytes, &key)
	file.Close()
	
	os.Setenv(common.EnvApiKeyID, key.Public)
	os.Setenv(common.EnvApiSecretKey, key.Secret)
	alpaca.SetBaseUrl("https://paper-api.alpaca.markets")

	/* Read stocks from file */
	file, err = os.Open("stocks.json")
	if err != nil {
		panic(err)
	}
	bytes, err = ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(bytes, &watchList)
	file.Close()
}

func runWebsite() {
	http.Handle("/", http.FileServer(
		http.Dir("./static"),
	))
	fmt.Println("Server Running")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}

func getClient() *alpaca.Client {
	/* Create Client */
	client := alpaca.NewClient(common.Credentials())
	_, err := client.GetAccount()
	if err != nil {
		panic(err)
	}
	return client
}

func quote(ticker string, client *alpaca.Client) {
	for {
		/* Get latest minute of data */
		barCount := 1
		bars, err := client.GetSymbolBars(ticker, alpaca.ListBarParams {
			Timeframe: "minute",
			Limit: &barCount,
		})
		if err != nil {
			panic(err)
		}

		/* Convert the bars to JSON */
		data, err := json.Marshal(bars)
		if err != nil {
			panic(err)
		}
		datastring := string(data)

		/* Create File */
		file, err := os.Create("static/stocks/" + ticker + ".json")
		if err != nil {
			panic(err)
		}

		/* Write data to file */
		_, err = io.WriteString(file, datastring)
		if err != nil {
			panic(err)
		}

		/* Close the file and wait */
		file.Close()
		t := rand.Intn(10) + 2
		fmt.Println(ticker + " " + datastring)
		time.Sleep(time.Duration(t) * time.Second)
	}
}

func main() {
	client := getClient()
	for _, ticker := range watchList.Stocks {
		go quote(ticker, client)
	}
	runWebsite()
}
