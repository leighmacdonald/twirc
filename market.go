package twirc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type MarketPrice struct {
	Success     bool `json:"success"`
	Name        string
	LowestPrice string `json:"lowest_price"`
	Volume      int    `json:"volume,string"`
	MedianPrice string `json:"median_price"`
}

func NormalizeItemName(name string) string {
	return strings.Title(name)
}

func GetPrice(hash_name string) (*MarketPrice, error) {
	name := NormalizeItemName(hash_name)
	url := fmt.Sprintf("http://steamcommunity.com/market/priceoverview/?currency=usd&appid=440&market_hash_name=%s", url.QueryEscape(name))
	var price MarketPrice
	err := loadJSON(url, &price)
	if err != nil {
		log.Println(err.Error())
	}
	price.Name = name
	return &price, err
}

func loadJSON(url string, output_struct interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return json.Unmarshal(body, output_struct)
}
