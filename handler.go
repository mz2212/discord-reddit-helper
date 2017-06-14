package main

import (
  "regexp"
  "strings"
  "fmt"

  "github.com/mz2212/discord_user_sim/markov"
  "github.com/bwmarrin/discordgo"
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
    client.ChannelMessageSend(message.ChannelID, "I don't work, and Tidest can't get this working...\nSorry.\nIf you've got an idea as to how to fix this, check me out on github.\nhttps://github.com/mz2212/discord-reddit-helper") //TODO: Figure this one out.
    return
  }
}

func userGenerate(loc string) string {
  harvest, err := redditbot.Listing("/u/" + loc, "")
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
/*
func subGenerate(loc string) string { // Sort of works...
  harvest, err := redditbot.Listing("/r/" + loc, "")
  if err != nil {
    fmt.Println("[Reddit] Failed to get listing for /r/", loc, ": ", err)
    msg := "Failed to get listing for " + loc + "\nEither that sub doesn't exist, or I bugged out..."
    return msg
  }
  gen := markov.New(2)
  for _, post := range harvest.Posts[:15] {
    for _, reply := range post.Replies[:3] { // Seems to get unhappy right about here.
      gen.Build(reply.Body)
    }
  }
  msg := gen.Generate(100) + " - /r/" + loc
  return msg
}
*/
