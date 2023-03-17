package handler

import (
	"log"
	"regexp"
	"strings"

	"github.com/tpc3/DuckPolice/lib/config"
	"github.com/tpc3/DuckPolice/lib/db"

	"github.com/bwmarrin/discordgo"
)

func urlCheck(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}()
	parsed := parseMsg(orgMsg.Content)
	for _, url := range parsed {
		found, channelid, messageid := db.SearchLog(orgMsg, &orgMsg.GuildID, &url)
		if found {
			if err := session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, config.CurrentConfig.Duplicate.React); err != nil {
				log.Print("Failed to add reaction: ", err)
				return
			}

			replyMessage := config.CurrentConfig.Duplicate.Message + "\nhttps://discord.com/channels/" + orgMsg.GuildID + "/" + channelid + "/" + messageid
			msg, err := session.ChannelMessageSendReply(orgMsg.ChannelID, replyMessage, orgMsg.Reference())
			if err != nil {
				log.Print("Failed to send message: ", err)
				return
			}

			if err := session.MessageReactionAdd(msg.ChannelID, msg.ID, config.CurrentConfig.Duplicate.Delete); err != nil {
				log.Print("Failed to add reaction: ", err)
				return
			}
		} else {
			db.AddLog(orgMsg, &orgMsg.GuildID, &url, &orgMsg.ChannelID, &orgMsg.ID)
		}
	}
}

var (
	re = regexp.MustCompile(`https?://[\w+.:?#[\]@!$&'()~*,;=/%-]+`)
)

func parseMsg(msg string) []string {
	if strings.HasPrefix(msg, "<") {
		return []string{}
	}

	var results []string
	for _, url := range re.FindAllString(msg, -1) {
		blacklistMatched := false
		for _, domain := range config.CurrentConfig.DomainBlacklist {
			if strings.Contains(url, domain) {
				blacklistMatched = true
				break
			}
		}
		if !blacklistMatched {
			var result string
			result = strings.ReplaceAll(url, "www.", "")
			result = strings.TrimSuffix(result, "/")
			result = strings.ReplaceAll(result, "youtu.be/", "youtube.com/watch?v=")
			result = strings.ReplaceAll(result, "m.", "")
			result = strings.ReplaceAll(result, "mobile.", "")
			results = append(results, result)
		}
	}
	return results
}
