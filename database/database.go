package database

import "github.com/keita0q/go-dsp/model"

type Database interface {
	LoadAllAdvertiser() ([]model.Advertiser, error)
	LoadAdvertiser(string) (*model.Advertiser, error)
	SaveAdvertiser(*model.Advertiser) error
}

type NotFoundError struct {
	message string
}

func NewNotFoundError(aMessage string) *NotFoundError {
	return &NotFoundError{message: aMessage}
}

func (aNotFoundError *NotFoundError) Error() string {
	return aNotFoundError.message
}