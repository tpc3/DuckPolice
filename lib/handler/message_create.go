package handler

import (
	"log"
	"strings"
	"time"

	"github.com/tpc3/DuckPolice/lib/cmds"
	"github.com/tpc3/DuckPolice/lib/config"
	"github.com/tpc3/DuckPolice/lib/db"

	"github.com/bwmarrin/discordgo"
)

func MessageCreate(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	guild := db.LoadGuild(&orgMsg.GuildID)
	prefix := guild.Prefix

	if config.CurrentConfig.Debug {
		start := time.Now()
		defer func() {
			log.Print("Message processed in ", time.Since(start).Milliseconds(), "ms.")
		}()
	}

	if orgMsg.Author.ID == session.State.User.ID {
		return
	}

	if orgMsg.Author.Bot {
		return
	}

	for _, v := range config.CurrentConfig.UserBlacklist {
		if orgMsg.Author.ID == v {
			return
		}
	}

	if !strings.HasPrefix(orgMsg.Content, prefix) {
		urlCheck(session, orgMsg)
		return
	}

	msg := strings.TrimSpace(orgMsg.Content[len(prefix):])
	if msg == "" {
		return
	}

	cmds.HandleCmd(session, orgMsg, guild, &msg)
}
