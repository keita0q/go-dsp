package model

type Bid struct {
	ID         string `json:"id"`
	FloorPrice float64 `json:"floor,omitempty"`
	Site       string `json:"site"`
	Page       string `json:"page"`
	Device     string `json:"device"`
	Browser    string `json:"browser"`
	Spot       string `json:"spot,omitempty"`
	Date       string `json:"date"`
	User       string `json:"user"`
	Test       int `json:"test"`
}

type Win struct {
	ID      string `json:"id"`
	Price   float64 `json:"price"`
	IsClick int `json:"isClick"`
}

type Advertiser struct {
	ID        string `json:"id"`
	Budget    int `json:"budget"`
	Cpc       int `json:"cpc"`
	NgDomains []string `json:"ngdomains"`
}