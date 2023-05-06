package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/DuckPolice/lib/config"
)

func ReactionCheck(session *discordgo.Session, msgreact *discordgo.MessageReactionAdd, msg *discordgo.Message, reactions []*discordgo.User) {
	reacted := false
	for _, user := range reactions {
		if user.ID != session.State.User.ID {
			reacted = true
		}
	}

	if reacted && (msgreact.Emoji.ID == config.CurrentConfig.Delete.Trigger || msgreact.Emoji.Name == config.CurrentConfig.Delete.Trigger) {
		session.MessageReactionRemove(msgreact.ChannelID, msg.MessageReference.MessageID, config.CurrentConfig.Duplicate.React, session.State.User.ID)
		session.ChannelMessageDelete(msgreact.ChannelID, msgreact.MessageID)
	}
}
