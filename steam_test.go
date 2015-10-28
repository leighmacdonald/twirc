package twirc

import (
	"io/ioutil"
	"log"
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

func TestSetGetSteamID(t *testing.T) {
	test_name := "rotorotorotorotorot"
	SqlDB.MustExec("DELETE FROM user WHERE username = ?", test_name)
	sid := SteamID("11111111111111111")
	err := SetSteamID(test_name, sid)
	if err != nil {
		t.Error(err)
	}
	sid2 := GetSteamID(test_name)
	log.Println(sid.ProfileURL())
	log.Println(sid2.ProfileURL())
	if sid2.ProfileURL() != sid.ProfileURL() {
		t.Error("Mismatched values")
	}
}
