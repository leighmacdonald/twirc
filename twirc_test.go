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
	inv, _ := FetchInventory("76561198084134025")
	if len(inv.Inventory) < 10 {
		t.Error("Did not decode enough entities")
	}
}

func TestResolveVanity(t *testing.T) {
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
	steam_id_1, e1 := NewSteamID("manofsnow")
	if e1 != nil {
		t.Error(e1.Error())
	}
	steam_id_2, e2 := NewSteamID("76561198005475714")
	if e2 != nil {
		t.Error(e1.Error())
	}
	if steam_id_1 != steam_id_2 {
		t.Error("Invalid steam id returned")
	}
}

func TestGetOrCreateUser(t *testing.T) {
	username := "test"
	DeleteUserByName(SqlDB, username)
	u1 := GetOrCreateUser(username)
	if u1.Username != username {
		t.Error("Invalid value")
	} else {
		u2 := GetOrCreateUser(username)
		if u2.UserID != u1.UserID {
			t.Error("Invalid value")
		}
	}
}
