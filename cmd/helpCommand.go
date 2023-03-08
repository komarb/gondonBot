package cmd

import (
	"Gondon/bot"
	"github.com/bwmarrin/discordgo"
)

func HelpCommand(ctx bot.Context) error {
	embed := bot.NewPlayerEmbed("GondonPlayer", "TRACK", "03:45", "@125")
	msgComplex := discordgo.MessageSend{Content: "TESTING", Embed: embed, Components: PlayerButtons}
	ctx.SendComplex(&msgComplex)
	return nil
}
