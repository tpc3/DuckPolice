package handler

import (
	"log"
	"runtime/debug"
	"strings"
	"time"

	"github.com/tpc3/DuckPolice/lib/cmds"
	"github.com/tpc3/DuckPolice/lib/config"
	"github.com/tpc3/DuckPolice/lib/db"

	"github.com/bwmarrin/discordgo"
)

func MessageCreate(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	defer func() {
		err := recover()
		if err != nil {
			log.Print("Oops, ", err)
			debug.PrintStack()
		}
	}()

	if config.CurrentConfig.Debug {
		start := time.Now()
		defer func() {
			log.Print("Message processed in ", time.Since(start).Milliseconds(), "ms.")
		}()
	}

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if orgMsg.Author.ID == session.State.User.ID || orgMsg.Content == "" {
		return
	}

	// Ignore all messages from blacklisted user
	for _, v := range config.CurrentConfig.UserBlacklist {
		if orgMsg.Author.ID == v {
			return
		}
	}

	for _, v := range config.CurrentConfig.ChannelBlacklist {
		if orgMsg.ChannelID == v {
			return
		}
	}

	// Ignore bot message
	if orgMsg.Author.Bot {
		return
	}

	guild := db.LoadGuild(&orgMsg.GuildID)

	isCmd := false
	var trimedMsg string
	if strings.HasPrefix(orgMsg.Content, guild.Prefix) {
		isCmd = true
		trimedMsg = strings.TrimPrefix(orgMsg.Content, guild.Prefix)
	} else if strings.HasPrefix(orgMsg.Content, session.State.User.Mention()) {
		isCmd = true
		trimedMsg = strings.TrimPrefix(orgMsg.Content, session.State.User.Mention())
		trimedMsg = strings.TrimPrefix(trimedMsg, " ")
	}
	if isCmd {
		if config.CurrentConfig.Debug {
			log.Print("Command processing")
		}
		cmds.HandleCmd(session, orgMsg, guild, &trimedMsg)
		return
	}

	urlCheck(session, orgMsg)
}
