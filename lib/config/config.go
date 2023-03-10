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
		React   string
		Message string
	}
	Db struct {
		Kind        string
		Parameter   string
		Tableprefix string
	}
	Guild           Guild
	DomainBlacklist []string `yaml:"url_blacklist"`
	UserBlacklist   []string `yaml:"user_blacklist"`
	LogPeriod       int64    `yaml:"log_period"`
}

type Guild struct {
	Prefix string
	Lang   string
}

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
	if CurrentConfig.LogPeriod == 0 {
		log.Print("LogPeriod is empty, setting 31536000")
	}

	loadLang()
}
