package handler

import (
	"github.com/tpc3/DuckPolice/lib/config"

	"github.com/bwmarrin/discordgo"
)

func MessageReactionAdd(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	msg, _ := session.ChannelMessage(reaction.ChannelID, reaction.MessageID)

	if reaction.UserID != session.State.User.ID && (reaction.Emoji.ID == config.CurrentConfig.Duplicate.Delete || reaction.Emoji.Name == config.CurrentConfig.Duplicate.Delete) {
		go session.MessageReactionRemove(reaction.ChannelID, msg.MessageReference.MessageID, config.CurrentConfig.Duplicate.React, session.State.User.ID)
		go session.ChannelMessageDelete(reaction.ChannelID, reaction.MessageID)
	}
}
