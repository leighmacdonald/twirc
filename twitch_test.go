package twirc

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestDecodeChatters(t *testing.T) {
	var chatters_resp ChatterResponse
	json_data, err := ioutil.ReadFile("./fixtures/chatters.json")
	if err != nil {
		t.Error(err.Error())
	}

	err = decodeChatters(json_data, &chatters_resp)
	if err != nil {
		t.Error(err.Error())
	}
	if chatters_resp.ChatterCount != 725 || chatters_resp.Count() != 725 {
		t.Errorf("Invalid cound: %d", chatters_resp.ChatterCount)
	}

}

func TestChatters(t *testing.T) {
	chatters, err := Chatters("twitch")
	if err != nil {
		t.Error(err.Error())
	}
	if chatters.ChatterCount <= 0 {
		t.Errorf("Invalid chatter count returned: %d", chatters.ChatterCount)
	}
}

func TestFetchEmotes(t *testing.T) {
	emotes, err := FetchEmotes()
	if err != nil {
		t.Error(err.Error())
	}
	if len(emotes.Emoticons) == 0 {
		t.Error("Invalid results returned")
	}
}
