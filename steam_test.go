package twirc

import (
	"io/ioutil"
	"testing"
)

func TestDecodePlayerSum(t *testing.T) {
	var player_sum ApiPlayerSumResponse
	json_data, err := ioutil.ReadFile("./fixtures/player_info.json")
	if err != nil {
		t.Error(err.Error())
	}

	err = decodePlayerSummary(json_data, &player_sum)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetPlayerSum(t *testing.T) {
	player, err := GetPlayerInfo(Conf.ApiKey, SteamID(Conf.SteamID))
	if err != nil {
		t.Error(err.Error())
	}
	if player.RealName != "Real Name" {
		t.Error("Invalid response decoded")
	}
}
