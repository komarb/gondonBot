package cmd

import "Gondon/bot"

func SkipCommand(ctx bot.Context) error {
	if !inVoiceChannel(&ctx) {
		return nil
	}
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if len(sess.Queue.Get()) == 0 {
		sess.Queue.Running = false
		bot.UpdatePlayerMsg(&ctx)
	}
	sess.Stop()
	return nil
}
