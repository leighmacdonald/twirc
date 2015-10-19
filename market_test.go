package twirc

import "testing"

func TestGetPrice(t *testing.T) {
	p, err := GetPrice("strange Rocket Launcher")
	if err != nil {
		t.Fatal(err.Error())
	}
	if p == nil || p.Volume <= 0 {
		t.Fatal("Invalid returned value")
	}
}
