package cmd

import "Gondon/bot"

func PopCommand(ctx bot.Context) error {
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if sess == nil {
		ctx.SendMsg("Not in a voice channel!")
		return nil
	}
	if sess.Queue.Pop() {
		ctx.SendMsg("Removed last song from queue.")
	} else {
		ctx.SendMsg("Queue is empty!")
	}
	return nil
}
