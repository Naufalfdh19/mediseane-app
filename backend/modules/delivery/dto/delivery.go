package dto

import "github.com/shopspring/decimal"

type (
	OngkirRequestDTO struct {
		OriginPostalCode      int `json:"origin_postal_code"`
		DestinationPostalCode int `json:"destination_postal_code"`
		Weight                int `json:"weight"`
	}

	UserCostResponseDTO struct {
		Cost int    `json:"cost"`
		Etd  string `json:"etd"`
	}

	OngkirResponseDTO struct {
		Id   int             `json:"id"`
		Name string          `json:"name,omitempty"`
		Cost decimal.Decimal `json:"cost,omitempty"`
		Etd  string          `json:"etd,omitempty"`
	}
	OngkirCostResponse struct {
		Data []struct {
			Name string          `json:"name"`
			Cost decimal.Decimal `json:"cost"`
			Etd  string          `json:"etd"`
		} `json:"data"`
	}
)
