package config

import (
	"errors"
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Debug   bool
	Help    string
	Discord struct {
		Token  string
		Status string
	}
	Db struct {
		Parameter   string
	}
	Guild     Guild
	LogPeriod int64 `yaml:"log_period"`
}

type Guild struct {
	Prefix       string
	Lang         string
	Ignore       []string
	ParsedIgnore []*regexp.Regexp `json:"-" yaml:"-"`
	Alert        struct {
		Type    string
		Message string
		React   string
		Reject  string
	}
	Domain struct {
		Mode string
		List []string
	}
	ChannelGroup map[string]string `yaml:"channel"`
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
	if CurrentConfig.LogPeriod == 0 {
		CurrentConfig.LogPeriod = 2592000
		log.Print("LogPeriod is empty, setting to 2592000")
	}

	loadLang()

	err = VerifyGuild(&CurrentConfig.Guild)
	if err != nil {
		log.Fatal("Config verify failed: ", err)
	}
}

func VerifyGuild(guild *Guild) error {
	if len(guild.Prefix) == 0 || len(guild.Prefix) >= 10 {
		return errors.New("prefix is too short or long")
	}
	_, exists := Lang[guild.Lang]
	if !exists {
		return errors.New("language does not exists")
	}
	return nil
}
