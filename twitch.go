package twirc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	httpClient = &http.Client{}
)

var (
	urlEmotes = "https://api.twitch.tv/kraken/chat/emoticons"
)

type (
	Chatter string

	ChatterResponse struct {
		ChatterCount int          `json:"chatter_count"`
		Chatters     ChatterTypes `json:"chatters"`
	}

	ChatterTypes struct {
		Moderators []Chatter `json:"moderators"`
		Staff      []Chatter `json:"staff"`
		Admins     []Chatter `json:"admins"`
		GlobalMods []Chatter `json:"global_mods"`
		Viewers    []Chatter `json:"viewers"`
	}
)

type (
	EmotesResp struct {
		links     map[string]string `json:"-"`
		Emoticons []EmotesContainer `json:"emoticons"`
	}

	EmotesContainer struct {
		Regex  string        `json:"regex"`
		Images []EmotesImage `json:"images"`
	}

	EmotesImage struct {
		EmoticonSet int    `json:"emoticon_set"`
		Height      int    `json:"height"`
		Width       int    `json:"width"`
		Url         string `json:"url"`
	}
)

func (c *ChatterResponse) Count() int {
	return len(c.Chatters.Admins) + len(c.Chatters.GlobalMods) + len(c.Chatters.Moderators) + len(c.Chatters.Staff) + len(c.Chatters.Viewers)
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

func FetchEmotes() (*EmotesResp, error) {
	resp_body, err := twitchAPIRequest(urlEmotes)
	var emotes EmotesResp
	err = json.Unmarshal(resp_body, &emotes)
	if err != nil {
		fmt.Println(err)
		return &emotes, err
	}
	return &emotes, nil
}

func twitchAPIRequest(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.twitchtv.v3+json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
