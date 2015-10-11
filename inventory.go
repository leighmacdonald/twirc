package twirc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type (
	EntryList []string

	Item struct {
		Id         string `json:"id"`
		ClassId    string `json:"classid"`
		InstanceID string `json:"instanceid"`
		Amount     string `json:"amount"`
		Pos        uint64 `json:"pos"`
	}

	MarketAction struct {
		Name string `json:"name"`
		Link string `json:"link"`
	}

	Action struct {
		Name string `json:"name"`
		Link string `json:"link"`
	}

	Tag struct {
		Internal_name string `json:"internal_name"`
		Name          string `json:"name"`
		Category      string `json:"category,omitempty"`
		Color         string `json:"color"`
		Category_name string `json:"category_name"`
	}

	Description struct {
		Type  string `json:"type,omitempty"`
		Value string `json:"value,omitempty"`
		Color string `json:"color,omitempty"`
	}

	AppData struct {
		Def_index string `json:"def_index"`
		Quality   string `json:"quality"`
	}

	Detail struct {
		Appid                         string         `json:"appid"`
		Classid                       string         `json:"classid"`
		Instanceid                    string         `json:"instanceid"`
		Icon_url                      string         `json:"icon_url"`
		Icon_url_large                string         `json:"icon_url_large"`
		Icon_drag_url                 string         `json:"icon_drag_url"`
		Name                          string         `json:"name"`
		Market_hash_name              string         `json:"market_hash_name"`
		Market_name                   string         `json:"market_name"`
		Name_colour                   string         `json:"name_color"`
		Background_color              string         `json:"background_color"`
		Type                          string         `json:"type"`
		Tradable                      uint64         `json:"tradable"`
		Marketable                    uint64         `json:"marketable"`
		Commodity                     uint64         `json:"commodity"`
		Market_tradable_restriction   string         `json:"market_tradable_restriction"`
		Market_marketable_restriction string         `json:"market_marketable_restriction"`
		Descriptions                  []Description  `json:"descriptions"`
		Actions                       []Action       `json:"actions"`
		MarketActions                 []MarketAction `json:"market_actions"`
		Tags                          []Tag          `json:"-"`
		AppData                       AppData        `json:"app_data"`
	}

	Inventory struct {
		Success      bool              `json:"success"`
		Inventory    map[string]Item   `json:"rgInventory"`
		Descriptions map[string]Detail `json:"rgDescriptions"`
		More         bool              `json:"more"`
		More_start   bool              `json:"more_start"`
	}
)

func FetchInventory(steam_id SteamID) (*Inventory, error) {
	var bits Inventory
	resp, err := http.Get(fmt.Sprintf("http://steamcommunity.com/profiles/%s/inventory/json/440/2/", steam_id))
	if err != nil {
		return &bits, err
	}
	defer resp.Body.Close()
	body, err_b := ioutil.ReadAll(resp.Body)
	if err_b != nil {
		return &bits, err_b
	}
	err_c := DecodeInventory(body, &bits)
	if err_c != nil {
		return &bits, err_c
	}
	return &bits, nil
}

func DecodeInventory(resp_body []byte, inv *Inventory) error {
	// The descriptions can be empty strings instead of empty arrays like we expect so we update
	// the json data if its the case so its able to be decoded properly.
	b2 := bytes.Replace(resp_body, []byte("\"descriptions\": \"\","), []byte("\"descriptions\": [],"), -1)
	err := json.Unmarshal(b2, inv)
	if err != nil {
		//fmt.Println(err)
		return err
	}
	return nil
}

func (inv *Inventory) FindMVMData() []MvMTour {

	keys := make([]string, 0)
	for k := range mvm_tags {
		keys = append(keys, k)
	}

	tours := make([]MvMTour, 0)

	for _, v := range inv.Descriptions {
		if stringInSlice(v.Name, keys) {
			tour_num_tmp := strings.Split(v.Type, " ")
			tour_count, conv_err := strconv.ParseUint(tour_num_tmp[1], 10, 64)
			if conv_err != nil {
				log.Println(conv_err.Error())
				continue
			}
			tour := NewMvMTour(v.Name, tour_count, mvm_tags[v.Name])
			for _, desc := range v.Descriptions {
				if stringInSlice(desc.Value, mvm_tags[v.Name]) {
					tour.AddCompleted(desc.Value)
				} else {
					fmt.Printf("not matched: %s", desc.Value)
				}
			}
			tours = append(tours, tour)
		}
	}

	return tours
}
