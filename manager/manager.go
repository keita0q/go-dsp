package manager

import (
	"github.com/keita0q/go-dsp/database"
	"github.com/keita0q/go-dsp/model"
	"github.com/keita0q/go-dsp/logic"
)

type Manager struct {
	logic       logic.Logic
	advertisers []model.Advertiser
}

type Config struct {
	Logic    logic.Logic
	Database database.Database
}

func New(aConfig *Config) (*Manager, error) {
	tAdvertisers, tError := aConfig.Database.LoadAllAdvertiser()
	if tError != nil {
		return nil, tError
	}
	return &Manager{logic: aConfig.Logic, advertisers: tAdvertisers}, nil
}

func (aManager *Manager) ExecuteCore(aBit *model.Bid) (*logic.Response, error) {
	return aManager.logic.Process(aBit, aManager.advertisers)
}
