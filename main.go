package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"

	"github.com/tpc3/DuckPolice/lib/config"
	"github.com/tpc3/DuckPolice/lib/db"
	"github.com/tpc3/DuckPolice/lib/handler"
)

func main() {
	token := config.CurrentConfig.Discord.Token
	discord, err := discordgo.New("Bot " + token)

	if err != nil {
		log.Fatalf("error creating Discord session: %s", err.Error())
	}

	discord.AddHandler(handler.MessageCreate)
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent | discordgo.IntentsGuildMessageReactions)

	err = discord.Open()

	if err != nil {
		log.Fatalf("error opening connection: %s", err.Error())
	}

	log.Printf("DuckPolice is now dispatching!")
	discord.UpdateGameStatus(0, config.CurrentConfig.Discord.Status)

	defer func() {
		discord.Close()
		db.Close()
		log.Print("DuckPolice is gracefully shutdowning!")
	}()

	go db.AutoLogCleaner()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
