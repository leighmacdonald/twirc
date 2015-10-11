package twirc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	vanity_url     = "http://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key=%s&vanityurl=%s"
	LastGameIP     = "*"
	UpdateGameData = true
)

type ApiVanityResponse struct {
	Response Vanity `json:"response"`
}

type Vanity struct {
	SteamID string `json:"steamid,omitempty"`
	Success int    `json:"success"`
	Message string `json:"message,omitempty"`
}

type SteamID string

type ApiPlayerSumResponse struct {
	Response Players `json:"response"`
}

type Players struct {
	Players []PlayerInfo `json:"players"`
}

type PlayerInfo struct {
	SteamID                  SteamID `json:"steamid"`
	CommunityVisibilityState int     `json:"communityvisibilitystate"`
	ProfileState             int     `json:"profilestate"`
	PersonaName              string  `json:"personaname"`
	LastLogoff               int     `json:"lastlogoff"`
	CommentPermission        int     `json:"commentpermission"`
	ProfileURL               string  `json:"profileurl"`
	Avatar                   string  `json:"avatar"`
	AvatarMedium             string  `json:"avatarmedium"`
	AvatarFull               string  `json:"avatarfull"`
	PersonaState             int     `json:"personastate"`
	RealName                 string  `json:"realname"`
	PrimaryClanID            string  `json:"primaryclanid"`
	TimeCreated              int     `json:"timecreated"`
	PersonastateFlags        int     `json:"personastateflags"`
	GameServerIP             string  `json:"gameserverip"`
	GameExtraInfo            string  `json:"gameextrainfo"`
	GameID                   string  `json:"gameid"`
}

func ResolveVanity(name string) (SteamID, error) {
	var api_resp ApiVanityResponse
	full_url := fmt.Sprintf(vanity_url, Conf.ApiKey, name)
	resp, err := http.Get(full_url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &api_resp)

	if err != nil {
		return "", err
	}
	return SteamID(api_resp.Response.SteamID), nil
}

func NewSteamID(steam_id string) (SteamID, error) {
	_, err := strconv.Atoi(steam_id)
	if err == nil {
		return SteamID(steam_id), nil
	}
	return ResolveVanity(steam_id)
}

func (sid SteamID) ProfileURL() string {
	return fmt.Sprintf("https://steamcommunity.com/profiles/%s", sid)
}

func SetSteamID(username string, steam_id SteamID) error {
	return db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte(DB_STEAM_ID)).Put([]byte(strings.ToLower(username)), []byte(steam_id))
		if err != nil {
			log.Println("Failed to save steam id")
			log.Println(err.Error())
		}
		return err
	})
}

func GetSteamID(username string) (SteamID, error) {
	var steam_id []byte
	err := db.View(func(tx *bolt.Tx) error {
		steam_id = tx.Bucket([]byte(DB_STEAM_ID)).Get([]byte(strings.ToLower(username)))
		if steam_id == nil {
			return errors.New("Unknown steam id")
		}
		return nil
	})
	return SteamID(steam_id), err
}

func decodePlayerSummary(resp_body []byte, player_sum *ApiPlayerSumResponse) error {
	err := json.Unmarshal(resp_body, player_sum)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetPlayerInfo(api_key string, steam_id SteamID) (*PlayerInfo, error) {
	var player_info PlayerInfo
	var info_response ApiPlayerSumResponse
	url := fmt.Sprintf(
		"https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s",
		api_key, steam_id)
	resp, err := http.Get(url)
	if err != nil {
		return &player_info, err
	}
	defer resp.Body.Close()
	body, err_b := ioutil.ReadAll(resp.Body)
	if err_b != nil {
		return &player_info, err_b
	}
	err_c := decodePlayerSummary(body, &info_response)
	if err_c != nil {
		return &player_info, err_c
	}

	if len(info_response.Response.Players) != 1 {
		return &player_info, errors.New("Invalid player count returned.")
	} else {
		return &info_response.Response.Players[0], nil
	}
}
