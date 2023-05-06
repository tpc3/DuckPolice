package common

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/DuckPolice/lib/config"
)

func getChannel(session *discordgo.Session, channelId string) *discordgo.Channel {
	channel, err := session.State.Channel(channelId)
	if err != nil {
		channel, _ = session.Channel(channelId)
	}
	return channel
}

func getParent(session *discordgo.Session, channelId string) string {
	if channelId == "" {
		return ""
	}
	channel := getChannel(session, channelId)
	return channel.ParentID
}

func GetGroup(session *discordgo.Session, guild *config.Guild, channelId string) string {

	if groupId, ok := guild.ChannelGroup[channelId]; ok {
		return groupId
	}

	ch := getChannel(session, channelId)
	if ch == nil {
		return ""
	}

	switch ch.Type {
	case discordgo.ChannelTypeGuildText, discordgo.ChannelTypeGuildVoice, discordgo.ChannelTypeGuildNews:
		if groupId, ok := guild.ChannelGroup["channel"]; ok {
			if groupId == "categoryId" {
				groupId = ch.ParentID
			} else if groupId == "channelId" {
				groupId = ch.ID
			}
			return groupId
		}
	case discordgo.ChannelTypeGuildNewsThread, discordgo.ChannelTypeGuildPublicThread, discordgo.ChannelTypeGuildPrivateThread:
		if groupId, ok := guild.ChannelGroup["thread"]; ok {
			if groupId == "categoryId" {
				groupId = getChannel(session, ch.ParentID).ParentID
			} else if groupId == "channelId" {
				groupId = ch.ParentID
			} else if groupId == "threadId" {
				groupId = ch.ID
			}
			return groupId
		}
	}

	parentId := getParent(session, channelId)

	for parentId != "" {
		if groupId, ok := guild.ChannelGroup[parentId]; ok {
			return groupId
		}
		parentId = getParent(session, parentId)
	}

	return "default"
}
