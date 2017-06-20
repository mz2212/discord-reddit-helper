package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mz2212/discord_user_sim/markov"
)

func messageCreate(client *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == client.State.User.ID { //Disregard messages from self.
		return
	}
	// '/r/' helper.
	matcher := regexp.MustCompile(`^[(reddit\.com)]?/?r/([^\s]+)`)
	if matcher.MatchString(message.Content) {
		matches := matcher.FindStringSubmatch(message.Content)
		//fmt.Println("Matches: ", matches) // Debugging.
		link := "Link: https://reddit.com/r/" + matches[1]
		client.ChannelMessageSend(message.ChannelID, link)
		return
	}

	// User/Sub quote generator.
	splitMessage := strings.Split(message.Content, " ")
	if splitMessage[0] == "#!usergen" {
		client.ChannelMessageSend(message.ChannelID, userGenerate(splitMessage[1]))
		return
	}
	if splitMessage[0] == "#!subgen" {
		client.ChannelMessageSend(message.ChannelID, subGenerate(splitMessage[1]))
		return
	}
}

func userGenerate(loc string) string {
	harvest, err := redditbot.Listing("/u/"+loc, "")
	if err != nil {
		fmt.Println("[Reddit] Failed to get listing for /u/", loc, ": ", err)
		msg := "Failed to get listing for " + loc + "\nEither that user doesn't exist, or I bugged out..."
		return msg
	}
	gen := markov.New(2)
	for _, comment := range harvest.Comments[:30] {
		gen.Build(comment.Body)
	}
	msg := gen.Generate(100) + " - /u/" + loc
	return msg
}

func subGenerate(loc string) string {
	harvest, err := redditbot.Listing("/r/"+loc+"/comments", "")
	if err != nil {
		fmt.Println("[Reddit] Failed to get listing for /r/", loc, ": ", err)
		msg := "Failed to get listing for " + loc + "\nEither that sub doesn't exist, or I bugged out..."
		return msg
	}
	gen := markov.New(2)
	for _, comment := range harvest.Comments[:30] {
		gen.Build(comment.Body)
	}
	msg := gen.Generate(100) + " - /r/" + loc
	return msg
}
