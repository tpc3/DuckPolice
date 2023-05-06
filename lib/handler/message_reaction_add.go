package handler

import (
	"log"
	"strings"

	"github.com/tpc3/DuckPolice/lib/db"

	"github.com/bwmarrin/discordgo"
)

func MessageReactionAdd(session *discordgo.Session, msgreact *discordgo.MessageReactionAdd) {
	guild := db.LoadGuild(msgreact.GuildID)

	if guild.Alert.Reject == "" {
		return
	}

	if msgreact.Emoji.ID != guild.Alert.Reject && msgreact.Emoji.Name != guild.Alert.Reject {
		return
	}

	msg, err := session.State.Message(msgreact.ChannelID, msgreact.MessageID)
	if err != nil {
		msg, err = session.ChannelMessage(msgreact.ChannelID, msgreact.MessageID)
		if err != nil {
			log.Panic("Failed to get channel message: ", err)
		}
	}

	if !strings.HasPrefix(msg.Content, guild.Alert.Message) {
		return
	}

	if msg.MessageReference != nil {
		session.MessageReactionRemove(msgreact.ChannelID, msg.MessageReference.MessageID, guild.Alert.React, session.State.User.ID)
	}
	session.ChannelMessageDelete(msgreact.ChannelID, msgreact.MessageID)
}
