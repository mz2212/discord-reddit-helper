package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"         // Probably a little overkill in this case, but ehh.
	"github.com/turnage/graw/reddit" // Seems like I don't need the functions from the main package.
)

var (
	redditbot reddit.Bot
)

func main() {
	viper.SetConfigName("bot")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("[Config] Error reading config: ", err)
		return
	}
	fmt.Println("[Discord] Initalizing...")
	// Begin Discord stuff...
	discord, err := discordgo.New("Bot " + viper.GetString("discord.key"))
	if err != nil {
		fmt.Println("[Discord] Error Creating discord session: ", err)
		return
	}
	discord.AddHandler(messageCreate)
	err = discord.Open()
	if err != nil {
		fmt.Println("[Discord] Error opening connection: ", err)
		return
	}
	// End Discord stuff...

	// Begin Reddit stuff...
	fmt.Println("[Reddit] Initalizing...")
	app := reddit.App{
		ID:       viper.GetString("reddit.id"),
		Secret:   viper.GetString("reddit.secret"),
		Username: viper.GetString("reddit.username"), // Why graw needs this, I don't know
		Password: viper.GetString("reddit.password"), // redd does just fine without...
	}
	botConfig := reddit.BotConfig{
		Agent: viper.GetString("reddit.agent"),
		App:   app,
		Rate:  1 * time.Second,
	}
	redditbot, err = reddit.NewBot(botConfig)
	if err != nil {
		fmt.Println("[Reddit] Error creating bot: ", err)
		return
	}
	// End Reddit Stuff...

	fmt.Println("[Main] Ready!")
	fmt.Println("[Main] Ctrl-C to Exit.")
	// Wait till Ctrl-C... Goddamn this is strange. Why not just loop it?
	sc := make(chan os.Signal, 1) // Channels to me, are magic. Perhaps once I learn threading, they won't be.
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc // This waits for somthing to come in on the "sc" channel.
	fmt.Println("[Main] Ctrl-C Recieved. Exiting!")
	discord.Close()
}
