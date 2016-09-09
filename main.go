package main

import (
	"encoding/json"
	"github.com/codegangsta/cli"
	"github.com/drone/routes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"github.com/keita0q/go-dsp/service"
	"github.com/keita0q/go-dsp/manager"
	"github.com/keita0q/go-dsp/logic/goLogic"
	"github.com/keita0q/go-dsp/database/local"
)

func main() {
	tApp := cli.NewApp()
	tApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "config, c",
		},
	}

	tApp.Action = func(aContext *cli.Context) {
		var tConfJSONPath string
		if aContext.String("config") != "" {
			tConfJSONPath = aContext.String("config")
		} else {
			tRunningProgramDirectory, tError := filepath.Abs(filepath.Dir(os.Args[0]))
			if tError != nil {
				log.Println("プログラムの走っているディレクトリの取得に失敗")
				os.Exit(1)
			}
			tConfJSONPath = path.Join(tRunningProgramDirectory, "config.json")
		}

		tJSONBytes, tError := ioutil.ReadFile(tConfJSONPath)
		if tError != nil {
			log.Println(tConfJSONPath + "の読み取りに失敗")
			os.Exit(1)
		}

		tConfig := &config{}
		if tError := json.Unmarshal(tJSONBytes, tConfig); tError != nil {
			log.Println(tError)
			log.Println(tConfJSONPath + "が不正なフォーマットです。")
			os.Exit(1)
		}

		tContextPath := "/" + tConfig.ContextPath + "/"

		// TODO implement Logic
		tLogic := goLogic.New()

		tDB := local.NewDatabase(tConfig.SavePath)

		tManager, tError := manager.New(&manager.Config{
			Logic: tLogic,
			Database: tDB,
		})
		if tError != nil {
			log.Println(tError)
			os.Exit(1)
		}

		tService := service.New(&service.Config{
			Manager: tManager,
			ContextPath:  tContextPath,
			ResourcePath: tConfig.ResourcePath,
		})
		if tError != nil {
			log.Println(tError)
			os.Exit(1)
		}

		tRouter := routes.New()

		tRouter.Post(path.Join(tContextPath, "/rest/v1/bid"), tService.BidRequest)
		tRouter.Post(path.Join(tContextPath, "/rest/v1/win"), tService.WinNotice)

		//tRouter.Post(path.Join(tContextPath, "/rest/v1/advertiser"), tService.SaveAdvertiser)
		tRouter.Get(path.Join(tContextPath, "/.*"), tService.GetFile)

		http.Handle(tContextPath, tRouter)

		http.ListenAndServe(":" + strconv.Itoa(tConfig.Port), nil)
	}

	tApp.Run(os.Args)
	os.Exit(0)
}

type config struct {
	ContextPath  string `json:"context_path"`
	Port         int    `json:"port"`
	SavePath     string `json:"save_path"`
	ResourcePath string `json:"resource_path"`
}

