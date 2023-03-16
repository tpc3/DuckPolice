package handler

import (
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
			go session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, config.CurrentConfig.Duplicate.React)
			replyMessage := config.CurrentConfig.Duplicate.Message + "\nhttps://discord.com/channels/" + orgMsg.GuildID + "/" + channelid + "/" + messageid
			msg, _ := session.ChannelMessageSendReply(orgMsg.ChannelID, replyMessage, orgMsg.Reference())
			go session.MessageReactionAdd(msg.ChannelID, msg.ID, config.CurrentConfig.Duplicate.Delete)
		} else {
			db.AddLog(orgMsg, &orgMsg.GuildID, &url, &orgMsg.ChannelID, &orgMsg.ID)
		}
	}
}

var (
	re = regexp.MustCompile(`https?://[\w+.:?#[\]@!$&'()~*,;=/%-]+`)
)

func parseMsg(origin string) []string {
	if strings.HasPrefix(origin, "<") {
		return []string{}
	}

	results := re.FindAllString(origin, -1)

	for i, result := range results {
		for _, domain := range config.CurrentConfig.DomainBlacklist {
			if strings.Contains(result, domain) {
				results = append(results[:i], results[i+1:]...)
				continue
			}
		}

		results[i] = strings.ReplaceAll(result, "www.", "")

		if strings.HasSuffix(result, "/") {
			results[i] = result[:len(result)-1]
		}

		results[i] = strings.ReplaceAll(results[i], "youtu.be/", "youtube.com/watch?v=")
		results[i] = strings.ReplaceAll(results[i], "m.", "")
		results[i] = strings.ReplaceAll(results[i], "mobile.", "")
	}

	return results
}
