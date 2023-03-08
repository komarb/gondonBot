package main

import (
	"Gondon/bot"
	"Gondon/cmd"
	"context"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"os/exec"
	"strings"
	"time"
)

var cfg *bot.Config
var DB *mongo.Client
var CmdHandler *bot.CommandHandler
var botId string
var Sessions *bot.SessionManager
var PlayerMessages *bot.PlayerMessages
var youtube *bot.Youtube
var PREFIX string
var textsColl *mongo.Collection
var imgsColl *mongo.Collection
var cursesColl *mongo.Collection

func init() {
	cfg = bot.GetConfigEnv()
	PREFIX = cfg.Prefix
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	initDB()
}

func initDB() {
	var err error
	log.WithField("dburi", cfg.DBUri).Info("Current database URI: ")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	DB, err = mongo.Connect(ctx, options.Client().ApplyURI(cfg.DBUri))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Can't connect to MongoDB, shutting down...")
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = DB.Ping(ctx, nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Failed to ping MongoDB, shutting down...")
	}
	dbName := cmd.GetDbName(os.Getenv("DB_URI"))
	log.WithFields(log.Fields{"db_name": dbName}).Info("Database information: ")
	textsColl = DB.Database("Gondon").Collection("texts")
	imgsColl = DB.Database("Gondon").Collection("images")
	cursesColl = DB.Database("Gondon").Collection("curses")
}

func main() {
	//checkForRequiredLibraries()
	CmdHandler = bot.NewCommandHandler()
	registerCommands()
	Sessions = bot.NewSessionManager()
	PlayerMessages = bot.NewPlayerMessages()
	youtube = &bot.Youtube{Cfg: cfg}
	discord, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Can't create Discord session, shutting down...")
	}
	usr, err := discord.User("@me")
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Can't get Discord account details, shutting down...")
	}
	botId = usr.ID
	// ! commands handler
	discord.AddHandler(textCommandsHandler)
	// game status handler
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		discord.UpdateGameStatus(0, cfg.GameStatus)
		guilds := discord.State.Guilds
		log.Info("Ready with ", len(guilds), " Discord guilds.")
	})
	// dj gondon player buttons handler and context menus
	discord.AddHandler(applicationCommandsHandler)
	//registerContextCommands(discord)
	//ResetCommands(discord)
	err = discord.Open()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Can't open Discord connection, shutting down...")
		return
	}
	log.Info("GandonBot has started!!!")

	defer discord.Close()
	<-make(chan struct{})
}

func registerContextCommands(discord *discordgo.Session) {
	for _, command := range cmd.MemCommands {
		_, err := discord.ApplicationCommandCreate("933147884818403458", "", command)
		//log.Info("registerContextCommands(discord)", " ", cmd.Name, " ", cmd.ID, " ", cmd.Type)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Fatal("Cannot create Application command")
		}
	}
}

func ResetCommands(discord *discordgo.Session) {
	cmds, _ := discord.ApplicationCommands("933147884818403458", "")
	for _, command := range cmds {
		err := discord.ApplicationCommandDelete("933147884818403458", "", command.ID)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Fatal("Cannot remove Application command")
		}
	}
}

func applicationCommandsHandler(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.User
	guild, err := discord.State.Guild(i.GuildID)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Can't get discord server, shutting down...")
	}
	channel, err := discord.State.Channel(i.ChannelID)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Can't get discord channel, shutting down...")
	}
	ctx := bot.NewContext(discord, guild, channel, user, &discordgo.MessageCreate{}, cfg, CmdHandler, Sessions, PlayerMessages, youtube, textsColl, imgsColl, cursesColl, bot.BullyingToday{})
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if handler, ok := cmd.MemCommandsHandlers[i.ApplicationCommandData().Name]; ok {
			handler(discord, i, ctx)
		}
		break
	case discordgo.InteractionMessageComponent:
		if handler, ok := cmd.PlayerButtonsHandlers[i.MessageComponentData().CustomID]; ok {
			handler(discord, i, ctx)
		}
		break
	default:
		break
	}

}
func textCommandsHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botId {
		return
	}
	channel, err := discord.State.Channel(message.ChannelID)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Can't get discord channel, shutting down...")
	}
	guild, err := discord.State.Guild(channel.GuildID)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Can't get discord guild, shutting down...")
	}
	ctx := bot.NewContext(discord, guild, channel, user, message, cfg, CmdHandler, Sessions, PlayerMessages, youtube, textsColl, imgsColl, cursesColl, bot.BullyingToday{})

	_, err = cmd.SaveMessageToDb(message.Message, ctx)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("Can't save message to DB")
	}
	content := message.Content
	if len(content) == 0 || content[:len(PREFIX)] != PREFIX {
		return
	}
	content = content[len(PREFIX):]
	args := strings.Fields(content)
	cmdName := strings.ToLower(args[0])
	log.Info(user.Username, message.GuildID, " : ", content)
	cmd, found := CmdHandler.Get(cmdName)
	if !found {
		return
	}
	ctx.Args = args[1:]
	command := *cmd
	err = command(*ctx)
	if err != nil {
		ctx.SendMsg("Something went wrong...")
		log.WithFields(log.Fields{"error": err}).Warn()
	}
}

func checkForRequiredLibraries() {
	//_, err := exec.LookPath("ffmpeg")
	//if err != nil {
	//	log.WithFields(log.Fields{"error": err}).Fatal("Required library (ffmpeg) is not installed")
	//}
	_, err := exec.LookPath("youtube-dl")
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Required library (youtube-dl) is not installed")
	}
}

//	func BullyInsert() {
//		data, err := os.Open("res/curses.txt")
//
//		if err != nil {
//			log.WithFields(log.Fields{"error": err}).Fatal("Can't read config.json file, shutting down...")
//		}
//		fileScanner := bufio.NewScanner(data)
//		fileScanner.Split(bufio.ScanLines)
//
//		for fileScanner.Scan() {
//
//			_, err := imgsColl.InsertOne(context.TODO(), interface{{Curse: "hello"}})
//
//		}
//
//		data.Close()
//	}
func registerCommands() {
	// Music
	CmdHandler.Register("help", cmd.HelpCommand, "Gives you this useful message!")
	CmdHandler.Register("play", cmd.PlayCommand, "Makes bot play music!")
	CmdHandler.Register("stop", cmd.StopCommand, "Makes bot stop music!")
	CmdHandler.Register("skip", cmd.SkipCommand, "Makes bot skip music!")
	CmdHandler.Register("pop", cmd.PopCommand, "Makes bot delete last song from queue!")
	CmdHandler.Register("bassboost", cmd.BassboostCommand, "Makes bot BASSING!")
	CmdHandler.Register("mod", cmd.ModCommand, "Makes bot play music different!")
	CmdHandler.Register("player", cmd.PlayerCommand, "Makes bot send you Player!")
	CmdHandler.Register("queue", cmd.QueueCommand, "Makes bot queue!")

	// Memes
	CmdHandler.Register("getallmsg", cmd.GetAllMsgCommand, "Gives you this help message!")
	CmdHandler.Register("mem", cmd.MemCommand, "Makes bot meme!")
	CmdHandler.Register("popusk", cmd.PopuskCommand, "Makes bot call someone a curseword!!")
	//Unused
	//CmdHandler.Register("add", cmd.AddCommand, "Makes bot add music!")
	//CmdHandler.Register("getallmsg", cmd.GetAllMsgCommand, "Gives you this help message!")
	//CmdHandler.Register("join", cmd.JoinCommand, "Makes bot join your channel!")
	//CmdHandler.Register("leave", cmd.LeaveCommand, "Makes bot leave your channel!")
	//CmdHandler.Register("yt", cmd.YoutubeCommand, "Makes bot search youtube videos!")
	//CmdHandler.Register("pick", cmd.PickCommand, "Makes bot pick from queue!")
	//CmdHandler.Register("play2", cmd.PlayOldCommand, "Makes bot play! (deprecated")

}
