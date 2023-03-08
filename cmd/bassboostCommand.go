package cmd

import "Gondon/bot"

func BassboostCommand(ctx bot.Context) error {
	if !inVoiceChannel(&ctx) {
		return nil
	}
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	sess.ToogleBassboost()
	return nil
}
