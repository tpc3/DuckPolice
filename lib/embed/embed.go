package embed

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	ColorURLChecker = 0xf4ebec
	ColorPink       = 0xf50057
)

var UnknownErrorNum int

func init() {
	UnknownErrorNum = 0
}

func NewEmbed(session *discordgo.Session, orgMsg *discordgo.MessageCreate) *discordgo.MessageEmbed {
	now := time.Now()
	msg := &discordgo.MessageEmbed{}
	msg.Author = &discordgo.MessageEmbedAuthor{}
	msg.Footer = &discordgo.MessageEmbedFooter{}
	msg.Author.IconURL = session.State.User.AvatarURL("256")
	msg.Author.Name = session.State.User.Username
	msg.Footer.IconURL = orgMsg.Author.AvatarURL("256")
	msg.Footer.Text = "Request from " + orgMsg.Author.Username + " @ " + now.String()
	return msg
}
