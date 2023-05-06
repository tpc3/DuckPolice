package config

import (
	"errors"
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
	Delete struct {
		Trigger string
	}
	Db struct {
		Kind        string
		Parameter   string
		Tableprefix string
	}
	Guild            Guild
	Domain           Domain
	UserBlacklist    []string `yaml:"user_blacklist"`
	ChannelBlacklist []string `yaml:"channel_blacklist"`
	LogPeriod        int64    `yaml:"log_period"`
}

type Guild struct {
	Prefix string
	Lang   string
}

type Domain struct {
	Type     string
	YamlList []string `yaml:"list"`
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

func SaveGuild(guild *Guild, domain *Domain) error {
	file, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return err
	}

	if guild.Prefix != CurrentConfig.Guild.Prefix || guild.Lang != CurrentConfig.Guild.Lang {
		CurrentConfig.Guild.Prefix = guild.Prefix
		CurrentConfig.Guild.Lang = guild.Lang
	}

	newConfig := Config{
		Debug:            config.Debug,
		Help:             config.Help,
		Discord:          config.Discord,
		Duplicate:        config.Duplicate,
		Delete:           config.Delete,
		Db:               config.Db,
		Guild:            *guild,
		Domain:           *domain,
		UserBlacklist:    config.UserBlacklist,
		ChannelBlacklist: config.ChannelBlacklist,
		LogPeriod:        config.LogPeriod,
	}
	data, err := yaml.Marshal(&newConfig)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFile, data, 0666)
	if err != nil {
		return err
	}

	CurrentConfig = newConfig

	return nil
}
