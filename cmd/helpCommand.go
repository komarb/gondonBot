package cmd

import (
	"Gondon/bot"
	"github.com/bwmarrin/discordgo"
)

func HelpCommand(ctx bot.Context) error {
	embed := bot.NewPlayerEmbed("GondonPlayer", "TEST", "01:23", ctx.User.Mention())
	msgComplex := discordgo.MessageSend{Content: "TESTING", Embed: embed, Components: PlayerButtons}
	ctx.SendComplex(&msgComplex)
	return nil
}
