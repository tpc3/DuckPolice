package cmds

import (
	"errors"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/DuckPolice/lib/common"
	"github.com/tpc3/DuckPolice/lib/config"
	"github.com/tpc3/DuckPolice/lib/db"
	"github.com/tpc3/DuckPolice/lib/embed"
)

const Config = "config"

func ConfigUsage(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, err error) {
	msg := embed.NewEmbed(session, orgMsg)
	if err != nil {
		msg.Title = config.Lang[guild.Lang].Error.Syntax
		msg.Description = err.Error() + "\n"
		msg.Color = embed.ColorPink
	}
	msg.Description += "`" + guild.Prefix + Config + " [<item> <value>]`\n" + config.Lang[guild.Lang].Usage.Config.Desc
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "prefix <prefix>",
		Value: config.Lang[guild.Lang].Usage.Config.Prefix,
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "lang <language>",
		Value: config.Lang[guild.Lang].Usage.Config.Lang,
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "alert type <type>\nalert message <message>\nalert react <emoji>\nalert reject <emoji>",
		Value: config.Lang[guild.Lang].Usage.Config.Alert,
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "domain mode <mode>\ndomain add <domain>\ndomain del <domain>",
		Value: config.Lang[guild.Lang].Usage.Config.Domain,
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "channel set [selector] <value>\nchannel del [selector]\nchannel test [channel]",
		Value: config.Lang[guild.Lang].Usage.Config.Channel,
	})
	common.ReplyEmbed(session, orgMsg, msg)
}

func ConfigCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message string) {
	split := strings.SplitN(message, " ", 2)
	if message == "" {
		msg := embed.NewEmbed(session, orgMsg)
		msg.Title = config.Lang[guild.Lang].CurrConf
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  "prefix",
			Value: guild.Prefix,
		})
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  "lang",
			Value: guild.Lang,
		})
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name: "alert",
			Value: "type: " + guild.Alert.Type +
				"\nmessage: " + guild.Alert.Message +
				"\nreact: " + guild.Alert.React +
				"\nreject: " + guild.Alert.Reject,
		})
		domainList := ""
		for _, v := range guild.Domain.List {
			domainList += v + "\n"
		}
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name: "domain",
			Value: "mode: " + guild.Domain.Mode +
				"\nlist: \n" + domainList,
		})
		channelList := ""
		if v, ok := guild.ChannelGroup["channel"]; ok {
			channelList += "channel: " + v + "\n"
		}
		if v, ok := guild.ChannelGroup["thread"]; ok {
			channelList += "thread: " + v + "\n"
		}
		for k, v := range guild.ChannelGroup {
			if k == "channel" || k == "thread" {
				continue
			}
			channelList += k + ": " + v + "\n"
		}
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  "channel",
			Value: channelList,
		})

		common.ReplyEmbed(session, orgMsg, msg)
		return
	}
	if len(split) != 2 {
		ConfigUsage(session, orgMsg, guild, errors.New("not enough arguments"))
		return
	}
	ok := false
	switch split[0] {
	case "prefix":
		guild.Prefix = split[1]
	case "lang":
		_, ok = config.Lang[split[1]]
		if ok {
			guild.Lang = split[1]
		} else {
			common.ErrorReply(session, orgMsg, "unsupported language")
			return
		}
	case "alert":
		alertCmd(session, orgMsg, guild, split[1])
	case "domain":
		domainCmd(session, orgMsg, guild, split[1])
	case "channel":
		channelCmd(session, orgMsg, guild, split[1])
	default:
		ConfigUsage(session, orgMsg, guild, errors.New("item not found"))
		return
	}
	err := db.SaveGuild(orgMsg.GuildID)
	if err != nil {
		common.UnknownError(session, orgMsg, guild.Lang, err)
		return
	}
	session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
}

var guildEmojiReactRegex = regexp.MustCompile(`<?a?:([a-zA-Z0-9_]+:\d+)>?`)
var guildEmojiIdRegex = regexp.MustCompile(`<?a?:[a-zA-Z0-9_]+:(\d+)>?`)

func alertCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message string) {
	split := strings.SplitN(message, " ", 2)
	if len(split) < 2 {
		ConfigUsage(session, orgMsg, guild, errors.New("not enough arguments"))
		return
	}
	switch split[0] {
	case "type":
		if split[1] != "reply" && split[1] != "message" && split[1] != "dm" {
			ConfigUsage(session, orgMsg, guild, errors.New("reply, msg, dm is allowed"))
			return
		}
		guild.Alert.Type = split[1]
	case "message":
		guild.Alert.Message = split[1]
	case "react":
		param := strings.TrimSpace(split[1])
		matched := guildEmojiReactRegex.FindStringSubmatch(param)
		emoji := ""
		if len(matched) == 2 {
			// custom emoji
			emoji = matched[1]
		} else if 1 < len(param) && len(param) < 8 && len([]rune(param)) == 1 {
			// Unicode emoji
			emoji = param
		} else {
			ConfigUsage(session, orgMsg, guild, errors.New("invalid emoji"))
			return
		}
		guild.Alert.React = emoji
	case "reject":
		param := strings.TrimSpace(split[1])
		matched := guildEmojiIdRegex.FindStringSubmatch(param)
		emoji := ""
		if len(matched) == 2 {
			// custom emoji
			emoji = matched[1]
		} else if 1 < len(param) && len(param) < 8 && len([]rune(param)) == 1 {
			// Unicode emoji
			emoji = param
		} else {
			ConfigUsage(session, orgMsg, guild, errors.New("invalid emoji"))
			return
		}
		guild.Alert.Reject = emoji
	default:
		ConfigUsage(session, orgMsg, guild, errors.New("unknown sub command"))
		return
	}
	err := db.SaveGuild(orgMsg.GuildID)
	if err != nil {
		common.UnknownError(session, orgMsg, guild.Lang, err)
		return
	}
	session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
}

func domainCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message string) {
	split := strings.SplitN(message, " ", 2)
	if len(split) < 2 {
		ConfigUsage(session, orgMsg, guild, errors.New("not enough arguments"))
		return
	}
	switch split[0] {
	case "mode":
		if split[1] != "white" && split[1] != "black" {
			ConfigUsage(session, orgMsg, guild, errors.New("white or black is allowed"))
			return
		}
		guild.Domain.Mode = split[1]
	case "add":
		guild.Domain.List = append(guild.Domain.List, split[1])
	case "del":
		for k, v := range guild.Domain.List {
			if split[1] == v {
				guild.Domain.List = append(guild.Domain.List[:k], guild.Domain.List[k+1:]...)
			}
		}
	default:
		ConfigUsage(session, orgMsg, guild, errors.New("unknown sub command"))
		return
	}
	err := db.SaveGuild(orgMsg.GuildID)
	if err != nil {
		common.UnknownError(session, orgMsg, guild.Lang, err)
		return
	}
	session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
}

func channelCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message string) {
	split := strings.SplitN(message, " ", 3)
	switch split[0] {
	case "set":
		if len(split) < 2 {
			ConfigUsage(session, orgMsg, guild, errors.New("not enough arguments"))
			return
		}
		selector := orgMsg.ChannelID
		value := split[1]
		if len(split) == 3 {
			selector = split[1]
			value = split[2]
		}
		guild.ChannelGroup[selector] = value
	case "del":
		selector := orgMsg.ChannelID
		if len(split) == 2 {
			selector = split[1]
		}
		delete(guild.ChannelGroup, selector)
	case "test":
		channel := orgMsg.ChannelID
		if len(split) == 2 {
			channel = split[1]
		}
		emb := embed.NewEmbed(session, orgMsg)
		emb.Title = "Group ID"
		emb.Description = common.GetGroup(session, guild, channel)
		common.ReplyEmbed(session, orgMsg, emb)
		return
	default:
		ConfigUsage(session, orgMsg, guild, errors.New("unknown sub command"))
		return
	}
	err := db.SaveGuild(orgMsg.GuildID)
	if err != nil {
		common.UnknownError(session, orgMsg, guild.Lang, err)
		return
	}
	session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
}
