package entity

import "github.com/shopspring/decimal"

type (
	OngkirLocation struct {
		Data []struct {
			ID int `json:"id"`
		} `json:"data"`
	}

	CalculateOngkir struct {
		BaseUrl       string
		SortingPrice  string
		OriginID      int
		DestinationID int
		Weight        int
		Courier       string
		Etd           int
	}

	OngkirRequest struct {
		OriginPostalCode      string
		DestinationPostalCode string
		Weight                int
		LogisticPartnerID     int
	}

	OngkirData struct {
		Id   int
		Name string
		Cost decimal.Decimal
		Etd  string
	}

	OngkirCostResponse struct {
		Data []struct {
			Name string          `json:"name"`
			Cost decimal.Decimal `json:"cost"`
			Etd  string          `json:"etd"`
		} `json:"data"`
	}

	UserCostResponse struct {
		Name string          `json:"name,omitempty" `
		Cost decimal.Decimal `json:"cost,omitempty"`
		Etd  string          `json:"etd,omitempty"`
	}
)
