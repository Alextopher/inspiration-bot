package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	bolt "go.etcd.io/bbolt"
)

// Get a link to a inspirational image
func getLink() (string, error) {
	resp, err := http.Get("https://inspirobot.me/api?generate=true")
	if err != nil {
		return "", err
	}

	// get the link out of the body
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	resp.Body.Close()

	return string(bytes), nil
}

type Inspiration struct {
	db      *bolt.DB
	session *discordgo.Session

	// schedule maps channel IDs to the hour of the day to send a message
	schedule map[string]int
}

func (inspiration *Inspiration) GetTargets(guild string) ([]string, error) {
	members, err := inspiration.session.GuildMembers(guild, "", 1000)
	if err != nil {
		return nil, err
	}

	targets := make([]string, len(members))
	for i, member := range members {
		targets[i] = member.User.Mention()
	}

	return targets, nil
}

// Schedule adds a job to the scheduler to send an inspirational message to a channel
// at a specific hour of the day (in UTC)
func (inspiration *Inspiration) Schedule(channel string, hour int) {
	inspiration.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("schedule"))
		h := fmt.Sprintf("%d", hour)
		b.Put([]byte(channel), []byte(h))
		return nil
	})
}

func (inspiration *Inspiration) Stop(channel string) {
	inspiration.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("schedule"))
		b.Delete([]byte(channel))
		return nil
	})
}

// Checks that the bot has access to all the channels in the schedule
func (inspiration *Inspiration) UpdateSchedule() {
	inspiration.db.Update(func(tx *bolt.Tx) error {
		count := 0

		b := tx.Bucket([]byte("schedule"))
		b.ForEach(func(k, _ []byte) error {
			channel := string(k)

			// Check if the bot has access to the channel
			_, err := inspiration.session.Channel(channel)
			if err != nil {
				b.Delete(k)
				count++
			}

			return nil
		})

		if count > 0 {
			log.Printf("Removed %d channels from the schedule\n", count)
		}

		return nil
	})
}

func (inspiration *Inspiration) RunScheduler() error {
	inspiration.UpdateSchedule()

	// Every hour on the hour check if we need to send an Inspiration message
	for {
		sleepUntilNextHour()

		// Map from channelID to hour of the day to send a message (in UTC)
		inspiration.schedule = make(map[string]int)

		// Get all the channels that have a scheduled Inspiration
		inspiration.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("schedule"))
			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				// Get the hour
				hour, err := strconv.Atoi(string(v))
				if err != nil {
					return err
				}

				// Get the channel
				channel := string(k)

				inspiration.schedule[channel] = hour
			}

			return nil
		})

		// Get the current hour
		hour := time.Now().UTC().Hour()

		// Check each channel
		for channelID, hourToSend := range inspiration.schedule {
			// If the hour matches, send the message
			if hour == hourToSend {
				err := inspiration.vibeCheck(nil, nil, channelID)
				if err != nil {
					log.Printf("Failed to vibecheck %s\nError: %s", channelID, err)
				}
			}
		}
	}
}

func (inspiration *Inspiration) vibeCheck(s *discordgo.Session, i *discordgo.InteractionCreate, channelID string) error {
	log.Println("Sending inspirational message to " + channelID)

	// Get a target
	channel, err := inspiration.session.Channel(channelID)
	if err != nil {
		return err
	}
	targets, err := inspiration.GetTargets(channel.GuildID)
	if err != nil {
		return err
	}
	target := targets[rand.Intn(len(targets))]

	// Get link
	link, err := getLink()
	if err != nil {
		return err
	}

	content := fmt.Sprintf("%s %s", link, target)

	if s == nil || i == nil {
		// Send the message to the channel
		_, err = inspiration.session.ChannelMessageSend(channelID, content)
		if err != nil {
			return err
		}
	} else {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func sleepUntilNextHour() {
	now := time.Now().UTC()
	next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, time.UTC)
	time.Sleep(next.Sub(now))
}
