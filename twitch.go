package twirc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Chatter string

type ChatterResponse struct {
	ChatterCount int          `json:"chatter_count"`
	Chatters     ChatterTypes `json:"chatters"`
}

type ChatterTypes struct {
	Moderators []Chatter `json:"moderators"`
	Staff      []Chatter `json:"staff"`
	Admins     []Chatter `json:"admins"`
	GlobalMods []Chatter `json:"global_mods"`
	Viewers    []Chatter `json:"viewers"`
}

func decodeChatters(resp_body []byte, chatters_resp *ChatterResponse) error {
	err := json.Unmarshal(resp_body, chatters_resp)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func Chatters(channel string) (*ChatterResponse, error) {
	var chatters_resp ChatterResponse
	url := fmt.Sprintf("https://tmi.twitch.tv/group/user/%s/chatters", channel)
	resp, err := http.Get(url)
	if err != nil {
		return &chatters_resp, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &chatters_resp, err
	}
	err = decodeChatters(body, &chatters_resp)
	if err != nil {
		return &chatters_resp, err
	}
	return &chatters_resp, nil
}
