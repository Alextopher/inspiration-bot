package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	bolt "go.etcd.io/bbolt"
)

var (
	inspirobot *Inspiration
)

func main() {
	rand.Seed(time.Now().UnixNano())

	log.Println("Starting bot...")
	godotenv.Load()

	// Load tokens from .env file.
	discordToken := os.Getenv("DISCORD_TOKEN")

	if discordToken == "" {
		log.Println("Please set DISCORD_TOKEN in the .env file.")
		return
	}

	// open the bolt key value store
	db, err := bolt.Open("./schedule.db", 0600, &bolt.Options{Timeout: time.Second})

	if err != nil {
		log.Println("Error opening bolt db: ", err)
		return
	}
	defer db.Close()

	// create the schedule bucket
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("schedule"))
		return nil
	})

	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Println("Error creating Discord session: ", err)
		return
	}

	inspirobot = &Inspiration{
		db:       db,
		session:  session,
		schedule: map[string]int{},
	}

	ch := make(chan struct{})
	session.AddHandler(func(_ *discordgo.Session, _ *discordgo.Ready) {
		log.Println("Bot is ready.")
		ch <- struct{}{}
	})

	session.Open()

	// number of scheduled channels
	scheduled := 0
	inspirobot.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("schedule"))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			scheduled++
		}

		return nil
	})
	log.Println("Schedule size: ", scheduled)

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
	inspirobot.RunScheduler()
}
