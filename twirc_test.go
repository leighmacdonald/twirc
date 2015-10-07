package twirc

import (
	"testing"
)

func TestFetchInventory(t *testing.T) {
	_, err := FetchInventory("76561198084134025")
	if err != nil {
		t.Error(err.Error())
	}
}
