package goLogic

import (
	"github.com/keita0q/go-dsp/logic"
	"github.com/keita0q/go-dsp/model"
	"time"
)

type GoLogic struct {

}

func New() *GoLogic {
	return &GoLogic{}
}

func (aLogic *GoLogic)Process(*model.Bid, []model.Advertiser) (*logic.Response, error) {

	// TODO implement

	time.Sleep(50 * time.Millisecond)
	return &logic.Response{ID:"hoge", BidPrice:6000, AdvertiserID:"adv1", Nurl:"htt://104.199.211.201:80/dsp/dsp/rest/v1/win"}, nil
}

