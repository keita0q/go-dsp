package model

// [WIP]

type Bid struct {
	ID string `json:"id"`
}

type Win struct {
	Price int `json:"price"`
}

type Advertiser struct {
	ID     string `json:"id"`
	Budget int `json:"budget"`
}