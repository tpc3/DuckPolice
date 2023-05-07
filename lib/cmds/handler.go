package cmds

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/tpc3/DuckPolice/lib/common"
	"github.com/tpc3/DuckPolice/lib/config"
)

func HandleCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message string) {
	splitMsg := strings.SplitN(message, " ", 2)
	var param string
	if len(splitMsg) == 2 {
		param = splitMsg[1]
	} else {
		param = ""
	}
	switch splitMsg[0] {
	case Ping:
		PingCmd(session, orgMsg, guild, param)
	case Help:
		HelpCmd(session, orgMsg, guild, param)
	case Config:
		ConfigCmd(session, orgMsg, guild, param)
	default:
		common.ErrorReply(session, orgMsg, config.Lang[guild.Lang].Error.NoCmd)
	}
}
