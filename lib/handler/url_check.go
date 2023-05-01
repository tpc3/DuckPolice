package handler

import (
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/tpc3/DuckPolice/lib/config"
	"github.com/tpc3/DuckPolice/lib/db"

	"github.com/bwmarrin/discordgo"
)

func urlCheck(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	parsed := parseMsg(orgMsg.Content)
	for _, url := range parsed {
		found, channelid, messageid := db.SearchLog(session, orgMsg, &orgMsg.GuildID, &url)
		message := config.CurrentConfig.Duplicate.Message + "\nhttps://discord.com/channels/" + orgMsg.GuildID + "/" + channelid + "/" + messageid
		if found {
			if config.CurrentConfig.Duplicate.DeleteMessage {
				session.ChannelMessageDelete(orgMsg.ChannelID, orgMsg.ID)
			}
			switch config.CurrentConfig.Duplicate.Alert {
			case "directmessage":
				dm, err := session.UserChannelCreate(orgMsg.Author.ID)
				if err != nil {
					log.Print("Create direct message channel error: ", err)
				}
				session.ChannelMessageSend(dm.ID, message)
			case "message":
				if !config.CurrentConfig.Duplicate.DeleteMessage {
					session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, config.CurrentConfig.Duplicate.React)
				}
				session.ChannelMessageSend(orgMsg.ChannelID, message)
			case "reply":
				if !config.CurrentConfig.Duplicate.DeleteMessage {
					session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, config.CurrentConfig.Duplicate.React)
					session.ChannelMessageSendReply(orgMsg.ChannelID, message, orgMsg.Reference())
				} else {
					session.ChannelMessageSend(orgMsg.ChannelID, message)
				}
			}
		} else {
			db.AddLog(orgMsg, &orgMsg.GuildID, &url, &orgMsg.ChannelID, &orgMsg.ID)
		}
	}
}

var (
	urlWithPathRegex = regexp.MustCompile(`https?://[\w+.:?#[\]@!$&'()~*,;=/%-]+`)
)

func parseMsg(origin string) []string {
	result := []string{}
	if strings.HasPrefix(origin, "< ") {
		return result
	}
url_loop:
	for _, rawUrl := range urlWithPathRegex.FindAllString(origin, -1) {
		rawUrl = strings.ToLower(rawUrl)
		rawUrl = strings.ReplaceAll(rawUrl, "youtu.be/", "youtube.com/watch?v=")
		Url, err := url.Parse(rawUrl)
		if err != nil {
			log.Panic("Invalid URL: ", rawUrl)
		}
		list := config.ListMap
		switch config.CurrentConfig.Domain.Type {
		case "black":
			if list[Url.Host] {
				continue url_loop
			}
		case "white":
			if !list[Url.Host] {
				continue url_loop
			}
		}
		Url.Host = strings.TrimPrefix(Url.Host, "www.")
		Url.Path = strings.TrimSuffix(Url.Path, "/")
		if Url.Host == "twitter.com" {
			Url.Host = strings.TrimPrefix(Url.Host, "m.")
			Url.Host = strings.TrimPrefix(Url.Host, "mobile.")
			Url.RawQuery = ""
		}
		result = append(result, Url.String())
	}
	return result
}
