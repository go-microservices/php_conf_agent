package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ClusterName string `yaml:"clusterName"`
	Type        int    `yaml:"type"`
	Address     string `yaml:"address"`
	Ip          string `yaml:"ip"`
	AutoIp      int    `yaml:"autoIp"`
	Configs     []struct {
		Path      string   `yaml:"path"`
		AppId     string   `yaml:"appId"`
		Namespace []string `yaml:"namespace"`
	} `yaml:"configs"`
}

var Conf *Config

var AppConfigPath string

func New() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	AppConfigPath = dir + "/app.yaml"
	configs, err := ioutil.ReadFile(AppConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(configs, &Conf)
	if err != nil {
		log.Fatal(err)
	}
}
