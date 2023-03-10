package cmds

import (
	"strconv"
	"time"

	"github.com/tpc3/DuckPolice/lib/config"
	"github.com/tpc3/DuckPolice/lib/db"
	"github.com/tpc3/DuckPolice/lib/embed"

	"github.com/bwmarrin/discordgo"
)

const Clean = "clean"

func CleanCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message *string) {
	result := embed.NewEmbed(session, orgMsg)
	result.Title = "Clean"
	start := time.Now()
	updated, err := db.CleanOldLog(&orgMsg.GuildID)
	if err != nil {
		UnknownError(session, orgMsg, &guild.Lang, err)
		return
	}
	if *updated != 0 {
		field := discordgo.MessageEmbedField{}
		field.Name = config.Lang[guild.Lang].Clean.Title
		field.Value = strconv.FormatInt(*updated, 10) + config.Lang[guild.Lang].Clean.DeletedLog
		result.Fields = append(result.Fields, &field)
	}
	result.Description += "clean old log: " + strconv.FormatInt(time.Since(start).Milliseconds(), 10) + "ms\n"
	ReplyEmbed(session, orgMsg, result)
}
