package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/tpc3/DuckPolice/lib/config"
	"github.com/tpc3/DuckPolice/lib/db"
	"github.com/tpc3/DuckPolice/lib/handler"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	Token := config.CurrentConfig.Discord.Token
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("error creating Discord session: ", err)
	}
	discord.AddHandler(handler.MessageCreate)
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent)
	err = discord.Open()
	if err != nil {
		log.Fatal("error opening connection: ", err)
	}
	discord.UpdateGameStatus(0, config.CurrentConfig.Discord.Status)
	go db.AutoLogCleaner()
	log.Print("DuckPolice is now dispatching!")
	defer discord.Close()
	defer db.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Print("DuckPolice is gracefully shutdowning!")
}
