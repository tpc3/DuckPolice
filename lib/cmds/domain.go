package cmds

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/DuckPolice/lib/config"
	"github.com/tpc3/DuckPolice/lib/db"
	"github.com/tpc3/DuckPolice/lib/embed"
)

const Domain = "domain"

func DomainCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message *string) {
	split := strings.SplitN(*message, " ", 2)
	if *message == "" {
		DomainUsage(session, orgMsg, guild, errors.New("not enough arguments"))
	}
	if (split[0] == "add" || split[0] == "del") && len(split) != 2 {
		DomainUsage(session, orgMsg, guild, errors.New("not enough arguments"))
		return
	}
	currentlist := config.CurrentConfig.Domain.YamlList
	currentmap := config.ListMap
	var list string
	var newlist []string
	switch split[0] {
	case "add":
		currentlist = append(currentlist, split[1])
		currentmap[split[1]] = true
		domain := config.Domain{
			Type:     config.CurrentConfig.Domain.Type,
			YamlList: currentlist,
		}
		err := config.SaveGuild(guild, &domain)
		if err != nil {
			UnknownError(session, orgMsg, &guild.Lang, err)
			return
		}
		session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
	case "del":
		newlist = RemoveElement(&currentlist, &split[1])
		currentmap[split[1]] = false
		domain := config.Domain{
			Type:     config.CurrentConfig.Domain.Type,
			YamlList: newlist,
		}
		err := config.SaveGuild(guild, &domain)
		if err != nil {
			UnknownError(session, orgMsg, &guild.Lang, err)
			return
		}
		session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
	case "list":
		for _, v := range currentlist {
			list += "- " + v + "\n"
		}
		embed := embed.NewEmbed(session, orgMsg)
		embed.Title = "Domain List"
		embed.Description = "```yaml\n" + list + "```"
		SendEmbed(session, orgMsg, embed)
	}
	err := db.SaveGuild(&orgMsg.GuildID, guild)
	if err != nil {
		UnknownError(session, orgMsg, &guild.Lang, err)
		return
	}
}

func RemoveElement(list *[]string, element *string) []string {
	currentlist := *list
	index := -1
	for i, v := range currentlist {
		if v == *element {
			index = i
			break
		}
	}
	if index == -1 {
		return currentlist
	}
	return append(currentlist[:index], currentlist[index+1:]...)
}

func DomainUsage(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, err error) {
	msg := embed.NewEmbed(session, orgMsg)
	if err != nil {
		msg.Title = config.Lang[guild.Lang].Error.Syntax
		msg.Description = err.Error() + "\n"
		msg.Color = embed.ColorPink
	}
	ReplyEmbed(session, orgMsg, msg)
}
