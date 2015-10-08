package twirc

import (
	"io/ioutil"
	"testing"
	"github.com/kr/pretty"
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

	//d := inv.FindMVMData()
	for _, v := range inv.Descriptions {
		pretty.Println(v)
	}

}

func TestFetchInventory(t *testing.T) {
	_, err := FetchInventory("76561198084134025")
	if err != nil {
		t.Error(err.Error())
	}
}
