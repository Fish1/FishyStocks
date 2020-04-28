package main

import (
	"os"
	"io"
	"io/ioutil"
	"fmt"
	"time"
	"encoding/json"
	"net/http"
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
)

type Key struct {
	Public string `json:"public"`
	Secret string `json:"secret"`
}

func init() {
	fmt.Println("Welcome to Fish Stock") 

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
		/* Get Stock Data */
		barCount := 5
		bars, err := client.GetSymbolBars(ticker, alpaca.ListBarParams {
			Timeframe: "day",
			Limit: &barCount,
		})
		if err != nil {
			panic(err)
		}

		/* Create File */
		file, err := os.Create("static/stocks/" + ticker + ".json")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Convert the bars to JSON
		data, err := json.Marshal(bars)
		if err != nil {
			panic(err)
		}
		datastring := string(data)

		/* Write data to file */
		_, err = io.WriteString(file, datastring)
		if err != nil {
			panic(err)
		}
		time.Sleep(30 * time.Second)
	}
}

func main() {
	client := getClient()
	go quote("AAPL", client)
	go quote("GOOG", client)
	runWebsite()
}
