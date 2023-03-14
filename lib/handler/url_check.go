package handler

import (
	"regexp"
	"strings"
	"time"

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
			session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
				if r.MessageID == msg.ID && r.UserID != s.State.User.ID {
					go s.MessageReactionRemove(msg.ChannelID, orgMsg.ID, config.CurrentConfig.Duplicate.React, s.State.User.ID)
					go s.MessageReactionsRemoveAll(msg.ChannelID, msg.ID)
					go s.ChannelMessageEdit(msg.ChannelID, msg.ID, config.CurrentConfig.Duplicate.Bye)
					go func() {
						<-time.After(3 * time.Second)
						s.ChannelMessageDelete(msg.ChannelID, msg.ID)
					}()
				}
			})

		} else {
			db.AddLog(orgMsg, &orgMsg.GuildID, &url, &orgMsg.ChannelID, &orgMsg.ID)
		}
	}
}

var (
	re = regexp.MustCompile(`https?://[\w+.:?#[\]@!$&'()~*,;=/%-]+`)
)

func parseMsg(origin string) []string {
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
