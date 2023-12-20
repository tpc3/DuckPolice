package handler

import (
	"log"
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/tpc3/DuckPolice/lib/common"
	"github.com/tpc3/DuckPolice/lib/config"
	"github.com/tpc3/DuckPolice/lib/db"

	"github.com/bwmarrin/discordgo"
)

func urlCheck(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	guild := db.LoadGuild(orgMsg.GuildID)

	groupId := common.GetGroup(session, guild, orgMsg.ChannelID)

	parsed := parseMsg(guild, orgMsg.Content)

	message := guild.Alert.Message

	for _, url := range parsed {
		found, msg := db.SearchLog(session, orgMsg.GuildID, groupId, &url)
		if found {
			message += "\nhttps://discord.com/channels/" + orgMsg.GuildID + "/" + msg.ChannelID + "/" + msg.MessageID
		} else {
			db.AddLog(orgMsg, orgMsg.GuildID, groupId, &url, orgMsg.ChannelID, orgMsg.ID)
		}
	}

	if message == guild.Alert.Message {
		return
	}

	switch guild.Alert.Type {
	case "dm":
		dm, err := session.UserChannelCreate(orgMsg.Author.ID)
		if err != nil {
			log.Print("Create direct message channel error: ", err)
		} else {
			session.ChannelMessageSend(dm.ID, message)
		}
	case "message":
		err := session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, guild.Alert.React)
		if err != nil {
			common.UnknownError(session, orgMsg, guild.Lang, err)
		}
		session.ChannelMessageSend(orgMsg.ChannelID, message)
	case "reply":
		err := session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, guild.Alert.React)
		if err != nil {
			common.UnknownError(session, orgMsg, guild.Lang, err)
		}
		reply := discordgo.MessageSend{}
		reply.Content = message
		reply.Reference = orgMsg.Reference()
		reply.AllowedMentions = &discordgo.MessageAllowedMentions{
			RepliedUser: false,
		}
		_, err = session.ChannelMessageSendComplex(orgMsg.ChannelID, &reply)
		if err != nil {
			log.Print("Failed to send reply: ", err)
		}
	}
}

var (
	urlWithPathRegex = regexp.MustCompile(`https?://[\w+.:?#[\]@!$&'()~*,;=/%-]+`)
)

func parseMsg(guild *config.Guild, origin string) []string {
	result := make(map[string]struct{})
	for _, v := range guild.ParsedIgnore {
		if v.MatchString(origin) {
			return []string{}
		}
	}
	for _, rawUrl := range urlWithPathRegex.FindAllString(origin, -1) {
		rawUrl = strings.ToLower(rawUrl)
		rawUrl = strings.ReplaceAll(rawUrl, "youtu.be/", "youtube.com/watch?v=")

		Url, err := url.Parse(rawUrl)
		if err != nil {
			log.Panic("Invalid URL: ", rawUrl)
		}

		Url.Host = strings.TrimPrefix(Url.Host, "www.")
		switch Url.Host {
		case "twitter.com", "m.twitter.com", "mobile.twitter.com", "fxtwitter.com", "vxtwitter.com", "x.com":
			Url.Host = "twitter.com"
			Url.RawQuery = ""
		}
		if net.ParseIP(Url.Host) != nil {
			continue
		}

		Url.Path = strings.TrimSuffix(Url.Path, "/")
		if Url.Path == "" {
			continue
		}

		found := false
		for _, v := range guild.Domain.List {
			if Url.Host == v {
				found = true
				break
			}
		}
		switch guild.Domain.Mode {
		case "black":
			if found {
				continue
			}
		case "white":
			if !found {
				continue
			}
		}

		result[Url.String()] = struct{}{}
	}

	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}

	return keys
}
