package models

type PriceDepth struct {
	Timestamp uint64 `json:"timestamp"`
	Asks      []Ask  `json:"asks"`
	Bids      []Bid  `json:"bids"`
}

type Ask struct {
	Price string `json:"price"`
}

type Bid struct {
	Price string `json:"price"`
}

type Price struct {
	Timestamp uint64  `db:"ts"`
	AskPrice  float64 `db:"ask_price"`
	BidPrice  float64 `db:"bid_price"`
}
