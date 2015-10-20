package twirc

import (
	"encoding/json"
	"fmt"
)

type (
	BPCurrencyResponse struct {
		Response BPCurrencyBase `json:"response"`
	}

	BPCurrencyBase struct {
		Success    int                            `json:"success"`
		Currencies map[string]BPCurrencyContainer `json:"currencies"`
		Name       string                         `json:"name"`
		URL        string                         `json:"url"`
	}

	BPCurrencyContainer struct {
		Quality    int        `json:"quality"`
		PriceIndex int        `json:"priceindex"`
		Single     string     `json:"single"`
		Plural     string     `json:"plural"`
		Round      int        `json:"round"`
		Blanket    int        `json:"blanket"`
		Craftable  string     `json:"craftable"`
		Tradable   string     `json:"tradable"`
		DefIndex   int        `json:"defindex"`
		Price      BPCurrency `json:"price"`
	}

	BPCurrency struct {
		Value      float64 `json:"value"`
		ValueHigh  float64 `json:"Value_high"`
		Currency   string  `json:"Currebct"`
		Difference float64 `json:""Difference`
	}
)

func decodeBPCurrencyResponse(resp_body []byte, resp *BPCurrencyResponse) error {
	err := json.Unmarshal(resp_body, resp)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
