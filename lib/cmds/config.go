package cmds

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
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
	ReplyEmbed(session, orgMsg, msg)
}

func ConfigCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message *string) {
	split := strings.SplitN(*message, " ", 2)
	if *message == "" {
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
		ReplyEmbed(session, orgMsg, msg)
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
			ErrorReply(session, orgMsg, "unsupported language")
			return
		}
	default:
		ConfigUsage(session, orgMsg, guild, errors.New("item not found"))
		return
	}
	err := db.SaveGuild(&orgMsg.GuildID, guild)
	if err != nil {
		UnknownError(session, orgMsg, &guild.Lang, err)
		return
	}
	session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
}
