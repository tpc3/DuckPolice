package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Debug   bool
	Help    string
	Discord struct {
		Token  string
		Status string
	}
	Duplicate struct {
		DeleteMessage bool
		Alert         string
		Message       string
		React         string
	}
	Db struct {
		Kind        string
		Parameter   string
		Tableprefix string
	}
	Guild  Guild
	Domain struct {
		Type     string
		YamlList []string `yaml:"list"`
	}
	UserBlacklist    []string `yaml:"user_blacklist"`
	ChannelBlacklist []string `yaml:"channel_blacklist"`
	LogPeriod        int64    `yaml:"log_period"`
}

type Guild struct {
	Prefix string
	Lang   string
}

var ListMap map[string]bool

const configFile = "./config.yml"

var CurrentConfig Config

func init() {
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal("Config load failed: ", err)
	}
	err = yaml.Unmarshal(file, &CurrentConfig)
	if err != nil {
		log.Fatal("Config parse failed: ", err)
	}

	//verify
	if CurrentConfig.Debug {
		log.Print("Debug is enabled")
	}
	if CurrentConfig.Discord.Token == "" {
		log.Fatal("Token is empty")
	}
	if CurrentConfig.Db.Tableprefix == "" {
		log.Fatal("Tableprefix is empty")
	}
	if CurrentConfig.Domain.Type != "white" && CurrentConfig.Domain.Type != "black" {
		log.Fatal("Domain type is invalid")
	}
	if CurrentConfig.LogPeriod == 0 {
		CurrentConfig.LogPeriod = 2592000
		log.Print("LogPeriod is empty, setting to 2592000")
	}

	ListMap = make(map[string]bool)
	for _, item := range CurrentConfig.Domain.YamlList {
		ListMap[item] = true
	}

	loadLang()
}
