package logic

import "github.com/keita0q/go-dsp/model"

type Logic interface {
	Process(*model.Bid, []model.Advertiser) (*Response, error)
}

type Response struct {
	ID string `json:"id"`
}