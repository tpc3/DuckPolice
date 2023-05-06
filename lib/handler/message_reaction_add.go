package handler

import (
	"log"

	"github.com/tpc3/DuckPolice/lib/config"

	"github.com/bwmarrin/discordgo"
)

func MessageReactionAdd(session *discordgo.Session, msgreact *discordgo.MessageReactionAdd) {
	msg, err := session.State.Message(msgreact.ChannelID, msgreact.MessageID)
	if err != nil {
		msg, err = session.ChannelMessage(msgreact.ChannelID, msgreact.MessageID)
		if err != nil {
			log.Panic("Failed to get channel message: ", err)
		}
	}

	if msgreact.Emoji.ID == "" && msgreact.Emoji.Name == "" {
		return
	}

	reactions, err := session.MessageReactions(msgreact.ChannelID, msgreact.MessageID, config.CurrentConfig.Delete.Trigger, 100, "", "")
	if err != nil {
		log.Panic("Failed to get msgreacts: ", err)
	}

	ReactionCheck(session, msgreact, msg, reactions)
}
