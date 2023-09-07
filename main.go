package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting bot...")
	godotenv.Load()

	// Load tokens from .env file.
	discordToken := os.Getenv("DISCORD_TOKEN")

	if discordToken == "" {
		log.Println("Please set DISCORD_TOKEN in the .env file.")
		return
	}

	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Println("Error creating Discord session: ", err)
		return
	}

	ch := make(chan struct{})
	session.AddHandler(func(_ *discordgo.Session, _ *discordgo.Ready) {
		log.Println("Bot is ready.")
		ch <- struct{}{}
	})

	session.Open()
	<-ch

	// Handle application commands
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := handlers[i.ApplicationCommandData().Name]; ok {
			log.Println(i.ApplicationCommandData().Name)
			handler(s, i)
		}
	})

	// Update the bot's interactions
	_, err = session.ApplicationCommandBulkOverwrite(session.State.User.ID, "", commands)
	if err != nil {
		log.Println("Error overwriting commands: ", err)
	}

	log.Println("Bot is running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
