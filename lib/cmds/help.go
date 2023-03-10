package cmds

import (
	"github.com/tpc3/DuckPolice/lib/config"
	"github.com/tpc3/DuckPolice/lib/embed"

	"github.com/bwmarrin/discordgo"
)

const Help = "help"

func HelpCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message *string) {
	msg := embed.NewEmbed(session, orgMsg)
	msg.Title = "Help"
	msg.Description = config.Lang[guild.Lang].Help + "\n" + config.CurrentConfig.Help
	ReplyEmbed(session, orgMsg, msg)
}
