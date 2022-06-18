package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"neko-bot/API"
)

func main() {
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/item/31582796.json?print=pretty")

	if err != nil {
		log.Fatalln(err)
	}

	var HNResponse API.HackerNewsContent
	json.NewDecoder(resp.Body).Decode(&HNResponse)
	fmt.Print(HNResponse)
}
