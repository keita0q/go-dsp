package service

import (
	"net/http"
	"strings"
	"fmt"
	"path/filepath"
	"net/url"
	"encoding/json"
	"log"
	"io/ioutil"
	"strconv"
	"github.com/keita0q/go-dsp/model"
	"github.com/keita0q/go-dsp/database"
	"github.com/keita0q/go-dsp/manager"
)

type Service struct {
	manager      *manager.Manager
	contextPath  string
	resourcePath string
}

type Config struct {
	Manager      *manager.Manager
	ContextPath  string
	ResourcePath string
}

func New(aConfig *Config) *Service {
	return &Service{
		manager: aConfig.Manager,
		contextPath:  aConfig.ContextPath,
		resourcePath: aConfig.ResourcePath,
	}
}

func (aService *Service) GetFile(aWriter http.ResponseWriter, aRequest *http.Request) {
	tPath := strings.TrimPrefix(aRequest.RequestURI, aService.contextPath)
	fmt.Println(tPath)
	if i := strings.Index(tPath, "?"); i > 0 {
		tPath = tPath[:i]
	}
	http.ServeFile(aWriter, aRequest, filepath.Join(aService.resourcePath, tPath))
}

func (aService *Service) BidRequest(aWriter http.ResponseWriter, aRequest *http.Request) {
	aService.handler(func(aQuerys url.Values, aRequestBody []byte) (int, interface{}, error) {
		tBid := &model.Bid{}

		if tError := json.Unmarshal(aRequestBody, tBid); tError != nil {
			return http.StatusBadRequest, nil, tError
		}
		_, tError := aService.manager.ExecuteCore(tBid)
		if tError != nil {
			return http.StatusInternalServerError, nil, tError
		}
		return http.StatusNoContent, nil, nil

	})(aWriter, aRequest)
}

func (aService *Service) WinNotice(aWriter http.ResponseWriter, aRequest *http.Request) {
	aService.handler(func(aQueries url.Values, aRequestBody []byte) (int, interface{}, error) {
		tWin := &model.Win{}
		if tError := json.Unmarshal(aRequestBody, tWin); tError != nil {
			return http.StatusBadRequest, nil, tError
		}
		//if tError := aService.manager.WinProcess(tWin); tError != nil {
		//	return http.StatusInternalServerError, nil, tError
		//}
		return http.StatusNoContent, nil, nil
	})(aWriter, aRequest)
}

func handleError(aError error) int {
	if _, ok := aError.(*database.NotFoundError); ok {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func (aService *Service) handler(aAPI func(url.Values, []byte) (int, interface{}, error)) func(http.ResponseWriter, *http.Request) {
	return func(aWriter http.ResponseWriter, aRequest *http.Request) {
		log.Printf("[INFO] access:%s", aRequest.RequestURI)
		defer aRequest.Body.Close()

		tResponseBody, tError := ioutil.ReadAll(aRequest.Body)
		if tError != nil {
			http.Error(aWriter, tError.Error(), http.StatusBadRequest)
		}
		tStatusCode, tResult, tError := aAPI(aRequest.URL.Query(), tResponseBody)
		if tError != nil {
			http.Error(aWriter, tError.Error(), tStatusCode)
			return
		}

		if tStatusCode == http.StatusNoContent {
			aWriter.WriteHeader(http.StatusNoContent)
			return
		}

		tBytes, tError := json.MarshalIndent(tResult, "", "  ")
		if tError != nil {
			http.Error(aWriter, tError.Error(), tStatusCode)
			return
		}

		aWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		aWriter.Header().Set("Content-Length", strconv.Itoa(len(tBytes)))
		aWriter.Header().Set("Access-Control-Allow-Origin", "*")
		aWriter.WriteHeader(tStatusCode)
		aWriter.Write(tBytes)
	}
}
