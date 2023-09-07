package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "inspiration",
		Description: "posts an inspirational quote from https://inspirobot.me/",
		Type:        discordgo.ChatApplicationCommand,
	},
	{
		Name:        "source",
		Description: "Link to the source code",
		Type:        discordgo.ChatApplicationCommand,
	},
}

var handlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
	"inspiration": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		link, err := getLink()
		if err != nil {
			sendError(s, i, fmt.Errorf("something went wrong ): "))
			log.Println(err)
		} else {
			sendMessage(s, i, link)
		}
	},
	"source": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		sendMessage(s, i, "https://github.com/Alextopher/inspiration-bot")
	},
}

func sendMessage(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})

	if err != nil {
		log.Println("Error responding to interaction: ", err)
	}
}

func sendError(s *discordgo.Session, i *discordgo.InteractionCreate, e error) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: e.Error(),
		},
	})

	if err != nil {
		log.Println("Error responding to interaction: ", err)
	}
}
