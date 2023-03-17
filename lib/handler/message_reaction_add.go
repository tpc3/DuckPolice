package handler

import (
	"github.com/tpc3/DuckPolice/lib/config"

	"github.com/bwmarrin/discordgo"
)

func MessageReactionAdd(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	msg, _ := session.State.Message(reaction.ChannelID, reaction.MessageID)
	if msg == nil {
		msg, _ = session.ChannelMessage(reaction.ChannelID, reaction.MessageID)
	}

	if reaction.Emoji.ID == "" && reaction.Emoji.Name == "" {
		return
	}

	reactions, _ := session.MessageReactions(reaction.ChannelID, reaction.MessageID, config.CurrentConfig.Duplicate.Delete, 100, "", "")

	reacted := false
	for _, r := range reactions {
		if r.ID == session.State.User.ID {
			reacted = true
		}
	}

	if reacted && reaction.UserID != session.State.User.ID && (reaction.Emoji.ID == config.CurrentConfig.Duplicate.Delete || reaction.Emoji.Name == config.CurrentConfig.Duplicate.Delete) {
		go session.MessageReactionRemove(reaction.ChannelID, msg.MessageReference.MessageID, config.CurrentConfig.Duplicate.React, session.State.User.ID)
		go session.ChannelMessageDelete(reaction.ChannelID, reaction.MessageID)
	}
}
