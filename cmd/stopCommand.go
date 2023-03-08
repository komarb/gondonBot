package cmd

import "Gondon/bot"

func StopCommand(ctx bot.Context) error {
	if !inVoiceChannel(&ctx) {
		return nil
	}
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if sess.Queue.HasNext() {
		sess.Queue.Clear()
	}
	sess.Queue.Running = false
	bot.UpdatePlayerMsg(&ctx)
	sess.Stop()
	ctx.Sessions.Leave(ctx.Discord, *sess)
	return nil
}
