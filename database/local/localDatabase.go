package local

import (
	"path/filepath"
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/keita0q/go-dsp/database"
	"github.com/keita0q/go-dsp/model"
	"path"
)

type LocalDB struct {
	savePath string
}

func NewDatabase(aSavePath string) *LocalDB {
	return &LocalDB{savePath: aSavePath}
}

func (aDatabase *LocalDB) LoadAllAdvertiser() ([]model.Advertiser, error) {
	tAdvertisers := []model.Advertiser{}

	//tError := aDatabase.eachFile(aDatabase.savePath, func(aFileInfo os.FileInfo, aBytes []byte) error {
	//	tAdvertiser := model.Advertiser{}
	//	if tError := json.Unmarshal(aBytes, &tAdvertiser); tError != nil {
	//		return tError
	//	}
	//	tAdvertisers = append(tAdvertisers, tAdvertiser)
	//	return nil
	//})
	tBytes, tError := ioutil.ReadFile(path.Join(aDatabase.savePath, "budgets.json"))
	if tError != nil {
		return nil, database.NewNotFoundError(tError.Error())
	}
	tBudgets := &database.Budgets{}
	if tError := json.Unmarshal(tBytes, tBudgets); tError != nil {
		return nil, tError
	}

	tBytes, tError = ioutil.ReadFile(path.Join(aDatabase.savePath, "ngdomains.json"))
	if tError != nil {
		return nil, database.NewNotFoundError(tError.Error())
	}
	tNgDomains := &database.NgDomains{}
	if tError := json.Unmarshal(tBytes, tNgDomains); tError != nil {
		return nil, tError
	}

	for tAdv, tBudget := range tBudgets.AdvBudget {
		tNgs := tNgDomains.AdvNgs[tAdv]
		tAdvertiser := model.Advertiser{ID:tAdv, Budget:tBudget.Budget, Cpc: tBudget.Cpc, NgDomains:tNgs}
		tAdvertisers = append(tAdvertisers, tAdvertiser)
	}
	return tAdvertisers, tError
}

func (aDatabase *LocalDB) LoadAdvertiser(aID string) (*model.Advertiser, error) {
	tFilePath := filepath.Join(aDatabase.savePath, aID)

	_, tError := os.Stat(tFilePath)
	if tError != nil {
		return nil, database.NewNotFoundError(tError.Error())
	}

	tByte, tError := ioutil.ReadFile(tFilePath)
	if tError != nil {
		return nil, tError
	}
	tAdvertiser := &model.Advertiser{}

	if tError := json.Unmarshal(tByte, tAdvertiser); tError != nil {
		return nil, tError
	}

	return tAdvertiser, nil
}

func (aDatabase *LocalDB) SaveAdvertiser(aAdvertiser *model.Advertiser) error {
	tPath := filepath.Join(aDatabase.savePath, aAdvertiser.ID)
	return aDatabase.saveObject(tPath, aAdvertiser)
}

func (aDatabase *LocalDB) RemoveAdvertiser(aID string) error {
	tFilePath := filepath.Join(aDatabase.savePath, aID)
	return aDatabase.removeObject(tFilePath)
}

func (aDatabase *LocalDB) eachFile(aPath string, aAction func(os.FileInfo, []byte) error) error {
	tFileInfos, tError := ioutil.ReadDir(aPath)
	if tError != nil {
		return database.NewNotFoundError(tError.Error())
	}
	for _, tFileInfo := range tFileInfos {
		if tFileInfo.IsDir() {
			continue
		}
		tByte, tError := ioutil.ReadFile(filepath.Join(aPath, tFileInfo.Name()))
		if tError != nil {
			return tError
		}
		if tError := aAction(tFileInfo, tByte); tError != nil {
			return tError
		}
	}
	return nil
}

func (aDatabase *LocalDB) saveObject(aPath string, aObject interface{}) error {
	if _, tError := os.Stat(aPath); tError != nil {
		if tError := os.MkdirAll(filepath.Dir(aPath), 0775); tError != nil {
			return tError
		}
	}

	tJSONBytes, tError := json.Marshal(aObject)
	if tError != nil {
		return tError
	}

	return ioutil.WriteFile(aPath, tJSONBytes, 0660)
}

func (aDatabase *LocalDB) removeObject(aPath string) error {
	_, tError := os.Stat(aPath)
	if tError != nil {
		return database.NewNotFoundError(tError.Error())
	}

	return os.Remove(aPath)
}