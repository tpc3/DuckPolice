package handler

import (
	"time"

	"github.com/tpc3/DuckPolice/lib/config"

	"github.com/bwmarrin/discordgo"
)

func MessageReactionAdd(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	msg, _ := session.ChannelMessage(reaction.ChannelID, reaction.MessageID)

	if reaction.UserID != session.State.User.ID && reaction.Emoji.Name == config.CurrentConfig.Duplicate.Delete {
		go session.MessageReactionRemove(reaction.ChannelID, msg.MessageReference.MessageID, config.CurrentConfig.Duplicate.React, session.State.User.ID)
		go session.MessageReactionsRemoveAll(reaction.ChannelID, reaction.MessageID)
		go session.ChannelMessageEdit(reaction.ChannelID, reaction.MessageID, config.CurrentConfig.Duplicate.Bye)
		go func() {
			<-time.After(3 * time.Second)
			session.ChannelMessageDelete(reaction.ChannelID, reaction.MessageID)
		}()
	}
}
