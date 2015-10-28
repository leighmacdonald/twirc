package twirc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	PersonaStateFlags        int     `json:"personastateflags"`
	GameServerIP             string  `json:"gameserverip"`
	GameExtraInfo            string  `json:"gameextrainfo"`
	GameID                   string  `json:"gameid"`
}

func ResolveVanity(name string) (SteamID, error) {
	var api_resp ApiVanityResponse
	full_url := fmt.Sprintf(vanity_url, Conf.ApiKey, name)
	err := loadJSON(full_url, &api_resp)
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

func (sid SteamID) MVMLobbyURL() string {
	return fmt.Sprintf("http://mvmlobby.com/profile/%s", sid)
}

func SetSteamID(username string, steam_id SteamID) error {
	id := GetSteamID(username)
	if id == "" {

	}
	tx := SqlDB.MustBegin()
	tx.MustExec("INSERT INTO user (username, steamid) VALUES (?, ?)", username, strings.ToLower(string(steam_id)))
	return tx.Commit()
}

func GetSteamID(username string) SteamID {
	var steamid string
	row := SqlDB.QueryRow("SELECT steamid FROM user WHERE username = ? LIMIT 1", strings.ToLower(username))
	row.Scan(&steamid)
	return SteamID(steamid)
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &player_info, err
	}
	err = decodePlayerSummary(body, &info_response)
	if err != nil {
		return &player_info, err
	}
	if len(info_response.Response.Players) != 1 {
		return &player_info, errors.New("Invalid player count returned.")
	} else {
		return &info_response.Response.Players[0], nil
	}
}
