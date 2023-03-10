package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Context struct {
	Discord        *discordgo.Session
	Guild          *discordgo.Guild
	VoiceChannel   *discordgo.Channel
	TextChannel    *discordgo.Channel
	User           *discordgo.User
	Message        *discordgo.MessageCreate
	Args           []string
	Conf           *Config
	CmdHandler     *CommandHandler
	Sessions       *SessionManager
	PlayerMessages *PlayerMessages
	Youtube        *Youtube
	TextsColl      *mongo.Collection
	ImgsColl       *mongo.Collection
	CursesColl     *mongo.Collection
	BullyingToday  BullyingToday
}

func NewContext(discord *discordgo.Session, guild *discordgo.Guild, textChannel *discordgo.Channel,
	user *discordgo.User, message *discordgo.MessageCreate, conf *Config, cmdHandler *CommandHandler,
	sessions *SessionManager, playerMessages *PlayerMessages, youtube *Youtube, textsColl *mongo.Collection, imgsColl *mongo.Collection, cursesColl *mongo.Collection) *Context {
	ctx := new(Context)
	ctx.Discord = discord
	ctx.Guild = guild
	ctx.TextChannel = textChannel
	ctx.User = user
	ctx.Message = message
	ctx.Conf = conf
	ctx.CmdHandler = cmdHandler
	ctx.Sessions = sessions
	ctx.PlayerMessages = playerMessages
	ctx.Youtube = youtube
	ctx.TextsColl = textsColl
	ctx.ImgsColl = imgsColl
	ctx.CursesColl = cursesColl
	return ctx
}

func (ctx *Context) SendMsg(msgText string) *discordgo.Message {
	msg, err := ctx.Discord.ChannelMessageSend(ctx.TextChannel.ID, msgText)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Can't send discord message")
		return nil
	}
	return msg
}

func (ctx *Context) SendComplex(data *discordgo.MessageSend) *discordgo.Message {
	msg, err := ctx.Discord.ChannelMessageSendComplex(ctx.TextChannel.ID, data)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Can't send discord complex")
		return nil
	}
	return msg
}

func (ctx *Context) GetVoiceChannel() *discordgo.Channel {
	if ctx.VoiceChannel != nil {
		return ctx.VoiceChannel
	}
	for _, state := range ctx.Guild.VoiceStates {
		if state.UserID == ctx.User.ID {
			channel, _ := ctx.Discord.State.Channel(state.ChannelID)
			ctx.VoiceChannel = channel
			return channel
		}
	}
	return nil
}
