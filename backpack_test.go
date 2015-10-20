package twirc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestDecodeBPCurrencyResponse(t *testing.T) {
	var resp BPCurrencyResponse
	json_data, err := ioutil.ReadFile("./fixtures/bp_tf_currencies.json")
	if err != nil {
		t.Error(err.Error())
	}

	err = decodeBPCurrencyResponse(json_data, &resp)
	if err != nil {
		t.Error(err.Error())
	}
	res2B, _ := json.Marshal(resp)
	fmt.Println(string(res2B))
}
