package twirc

import (
	"io/ioutil"
	"testing"
)

func TestDecodeInventory(t *testing.T) {
	var inv Inventory
	json_data, err := ioutil.ReadFile("./fixtures/inventory.json")
	if err != nil {
		t.Error(err.Error())
	}

	err = DecodeInventory(json_data, &inv)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestFetchInventory(t *testing.T) {
	_, err := FetchInventory("76561198084134025")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestResolveVanity(t *testing.T) {
	LoadConfig()
	steam_id, err := ResolveVanity("manofsnow")
	if err != nil {
		t.Log(err.Error())
		t.Error("Got error!")
	} else {
		if steam_id != "76561198005475714" {
			t.Error("Invalid response from API!")
		}
	}
}

func TestSteamID(t *testing.T) {
	steam_id_1 := SteamID("manofsnow")
	steam_id_2 := SteamID("76561198005475714")
	if steam_id_1 != steam_id_2 {
		t.Error("Invalid steam id returned")
	}

}
