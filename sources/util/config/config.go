package config

import (
	"fmt"
	"github.com/sdvdxl/go-tools/collections"
	"github.com/sdvdxl/wine/sources/bean"
	. "github.com/sdvdxl/wine/sources/util/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var (
	AuthPageConfig bean.AuthPageConfig
)

type authPageConfig struct {
	LoginPages []string `yaml:"login_pages"`
	LoginPaths []string `yaml:"login_paths"`
}

func init() {
	Logger.Info("init config ...")
	fmt.Println("init config")
	Logger.Info("read auth_pages.yaml")
	file, err := os.Open("properties/auth.pages.yaml")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var pageConfig authPageConfig
	err = yaml.Unmarshal(bytes, &pageConfig)
	if err != nil {
		panic(err)
	}

	pages := collections.NewSet(10)
	for _, v := range pageConfig.LoginPaths {
		pages.Add(v)
	}
	AuthPageConfig.LoginPaths = pages

	pages = collections.NewSet(10)
	for _, v := range pageConfig.LoginPages {
		pages.Add(v)
	}
	AuthPageConfig.LoginPages = pages

	Logger.Info("required login pages: %v", AuthPageConfig.LoginPages)
	Logger.Info("login paths: %v", AuthPageConfig.LoginPaths)
	Logger.Info("config inited")
}
