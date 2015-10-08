package twirc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	vanity_url = "http://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key=%s&vanityurl=%s"
)

type ApiVanityResponse struct {
	Response Vanity `json:"response"`
}

type Vanity struct {
	SteamID string `json:"steamid,omitempty"`
	Success int    `json:"success"`
	Message string `json:"message,omitempty"`
}

func ResolveVanity(name string) (string, error) {
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
	return api_resp.Response.SteamID, nil
}

func SteamID(steam_id string) (string, error) {
	out_str := ""
	if len(steam_id) == 17 {
		_, err := strconv.Atoi(steam_id)
		if err == nil {
			return steam_id, nil
		}
	}
	out_str, err := ResolveVanity(steam_id)
	if err != nil {
		return out_str, err
	}
	return out_str, nil
}
