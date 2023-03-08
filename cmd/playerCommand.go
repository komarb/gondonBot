package cmd

import (
	"Gondon/bot"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func PlayerCommand(ctx bot.Context) error {
	pm := ctx.PlayerMessages.Get(ctx.Guild.ID)
	if pm == nil {
		embed := bot.NewPlayerEmbed("🔇Stopped", "-", "-:-", "-")
		msgComplex := discordgo.MessageSend{Embed: embed, Components: PlayerButtons}
		newPlayerMessage := ctx.SendComplex(&msgComplex)
		ctx.PlayerMessages.Update(ctx.Guild.ID, newPlayerMessage)
		return nil
	}
	oldMsg, err := ctx.Discord.ChannelMessage(pm.ChannelID, pm.ID)
	if err != nil {
		return fmt.Errorf("can't get discord message: %s", err)
	}
	newMsg := ctx.SendComplex(&discordgo.MessageSend{Components: oldMsg.Components, Embeds: oldMsg.Embeds})
	err = ctx.Discord.ChannelMessageDelete(oldMsg.ChannelID, oldMsg.ID)
	if err != nil {
		return fmt.Errorf("can't delete discord message: %s", err)
	}
	ctx.PlayerMessages.Update(ctx.Guild.ID, newMsg)
	return nil
}
