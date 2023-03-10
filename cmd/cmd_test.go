package cmd

import (
	"Gondon/bot"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"testing"
)

var cfg *bot.Config
var ctx *bot.Context
var discord *discordgo.Session
var Sessions *bot.SessionManager
var PlayerMessages *bot.PlayerMessages
var youtube *bot.Youtube

func TestMain(m *testing.M) {
	ch := make(chan struct{})
	go setup(ch)
	<-ch
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup(ch chan struct{}) {
	var err error
	cfg = bot.GetConfigEnv()
	discord, err = discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		log.Fatal("Can't create Discord session, shutting down...")
	}
	err = discord.Open()
	if err != nil {
		log.Fatal("Can't open Discord connection, shutting down...")
	}

	user, _ := discord.User("333943790802960385")                //TODO: config for test
	channel, err := discord.State.Channel("1027663560706375731") //TODO: config for test
	guild, err := discord.State.Guild(channel.GuildID)
	ctx = bot.NewContext(discord, guild, channel, user, nil, cfg, nil, Sessions, PlayerMessages, youtube, nil, nil, nil)
	ch <- struct{}{}
	<-make(chan struct{})
}
func teardown() {
	discord.Close()
}

func TestHelpCommand(t *testing.T) {
	replyMsg, err := HelpCommand(*ctx)
	if err != nil || replyMsg.Content != "TESTING" {
		t.Errorf("Expected: <%s>, got: <%s>", "TESTING", replyMsg.Content)
	}
}
