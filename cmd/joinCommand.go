package cmd

import (
	"Gondon/bot"
	"fmt"
)

func JoinCommand(ctx bot.Context) error {
	if ctx.Sessions.GetByGuild(ctx.Guild.ID) != nil {
		return nil
	}
	vc := ctx.GetVoiceChannel()
	if vc == nil {
		ctx.SendMsg("You must be in a voice channel to use dj Gandon!")
		return fmt.Errorf("user not in vc")
	}
	_, err := ctx.Sessions.Join(ctx.Discord, ctx.Guild.ID, vc.ID, bot.JoinProperties{
		Muted:    false,
		Deafened: true,
	})
	if err != nil {
		ctx.SendMsg("An error occured!")
		return fmt.Errorf("can't join vc: %s", err)
	}
	return nil
}
