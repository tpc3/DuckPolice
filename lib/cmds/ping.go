package cmds

import (
	"github.com/tpc3/DuckPolice/lib/config"
	"github.com/tpc3/DuckPolice/lib/embed"

	"github.com/bwmarrin/discordgo"
)

const Ping = "ping"

func PingCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message *string) {
	embedMsg := embed.NewEmbed(session, orgMsg)
	embedMsg.Title = "Pong!"
	ReplyEmbed(session, orgMsg, embedMsg)
}
