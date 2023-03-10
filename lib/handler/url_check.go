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
		found, channelid, messageid := db.SearchLog(orgMsg, &orgMsg.GuildID, &url)
		if found {
			session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, config.CurrentConfig.Duplicate.React)
			session.ChannelMessageSendReply(orgMsg.ChannelID, config.CurrentConfig.Duplicate.Message+"\nhttps://discord.com/channels/"+orgMsg.GuildID+"/"+channelid+"/"+messageid, orgMsg.Reference())
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
		for _, v := range config.CurrentConfig.DomainBlacklist {
			if strings.Contains(Url.Host, v) {
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
