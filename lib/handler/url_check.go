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
	parsed := parseMsg(orgMsg.Content)
	log.Print(parsed)
	for _, url := range parsed {
		found, channelid, messageid := db.SearchLog(orgMsg, &orgMsg.GuildID, &url)
		if found {
			go func() {
				session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, config.CurrentConfig.Duplicate.React)
			}()
			replyMessage := config.CurrentConfig.Duplicate.Message + "\nhttps://discord.com/channels/" + orgMsg.GuildID + "/" + channelid + "/" + messageid
			msg, _ := session.ChannelMessageSendReply(orgMsg.ChannelID, replyMessage, orgMsg.Reference())
			session.MessageReactionAdd(msg.ChannelID, msg.ID, "kaere:1085113264402354186")
		} else {
			db.AddLog(orgMsg, &orgMsg.GuildID, &url, &orgMsg.ChannelID, &orgMsg.ID)
		}
	}
}

var (
	urlWithPathRegex = `((http|https):\/\/)?[\w\-]+(\.[\w\-]+)+[/\w\-\.\?\,\'\/\\\+&amp;%\$#\=~]*`
)

func parseMsg(origin string) []string {
	re := regexp.MustCompile(urlWithPathRegex)
	results := re.FindAllString(origin, -1)

	for i, result := range results {
		results[i] = strings.ReplaceAll(result, "www.", "")

		if strings.HasSuffix(result, "/") {
			results[i] = result[:len(result)-1]
		}

		if strings.Contains(result, "youtu.be") || strings.Contains(result, ".mobile") {
			results[i] = strings.ReplaceAll(results[i], "youtu.be/", "youtube.com/watch?v=")
			results[i] = strings.ReplaceAll(results[i], ".mobile", "")
		} else if strings.Contains(result, "twitter.com") && (strings.Contains(result, "m.") || strings.Contains(result, "mobile.")) {
			results[i] = strings.ReplaceAll(results[i], "m.", "")
			results[i] = strings.ReplaceAll(results[i], "mobile.", "")
		}
	}

	return results
}
