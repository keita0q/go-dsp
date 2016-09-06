package goLogic

import (
	"github.com/keita0q/go-dsp/logic"
	"github.com/keita0q/go-dsp/model"
)

type GoLogic struct {

}

func New() *GoLogic {
	return &GoLogic{}
}

func (aLogic *GoLogic)Process(*model.Bid, []model.Advertiser) (*logic.Response, error) {

	// TODO implement

	return nil, nil
}

