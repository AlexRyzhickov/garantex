package models

type PriceDepth struct {
	Timestamp uint64 `json:"timestamp"`
	Asks      []Ask  `json:"asks"`
	Bids      []Bid  `json:"bids"`
}

type Ask struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}

type Bid struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}

type Price struct {
	Timestamp uint64  `db:"ts"`
	AskPrice  float64 `db:"ask_price"`
	BidPrice  float64 `db:"bid_price"`
}
