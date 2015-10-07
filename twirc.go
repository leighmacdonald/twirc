package twirc

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/thoj/go-ircevent"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	Conf Config
)

type (
	Config struct {
		Server   string
		Name     string
		Password string
	}

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
		Value string `json:"value"`
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

func NewIrcClient(config Config) (*irc.Connection, error) {
	irc_conn := irc.IRC(config.Name, config.Name)
	irc_conn.VerboseCallbackHandler = true
	irc_conn.Debug = true
	irc_conn.Password = config.Password

	irc_conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		log.Info(e.Message())
		if e.Message() == "Test Message" {
			irc_conn.Quit()
		}
		if strings.HasPrefix("~inv", e.Message()) {
			log.Println("Got ~inv request")
			_, err := FetchInventory("76561198084134025")
			if err != nil {
				irc_conn.Privmsg("#roto_", "got inv")
			} else {
				irc_conn.Privmsg("#roto_", err.Error())
			}
		}
	})

	irc_conn.AddCallback("001", func(e *irc.Event) {
		irc_conn.Join("#manofsnow")
		irc_conn.Join("#cuddli")
		irc_conn.Join("#roto_")
	})

	err := irc_conn.Connect(config.Server)
	if err != nil {
		return nil, err
	}

	return irc_conn, nil
}

func FetchInventory(steam_id string) (Inventory, error) {
	var bits Inventory

	resp, err := http.Get(fmt.Sprintf("http://steamcommunity.com/profiles/%s/inventory/json/440/2/", steam_id))
	if err != nil {
		return bits, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &bits)
	if err != nil {
		return bits, err
	}

	b, _ := json.MarshalIndent(bits, "", "  ")
	os.Stdout.Write(b)

	return bits, nil
}

func LoadConfig() {
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		// handle error
		log.Fatalln(err.Error())
	}
}
