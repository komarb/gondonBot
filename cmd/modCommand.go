package cmd

import "Gondon/bot"

var Mods = map[string]string{
	"default": "asetrate=48000",
	"spedup":  "asetrate=48000*1.25,aresample=48000",
	"slowed":  "asetrate=48000*0.75,aresample=48000,aecho=1.0:0.8:50:0.5",
}

func ModCommand(ctx bot.Context) error {
	if !inVoiceChannel(&ctx) {
		return nil
	}
	argsLen := len(ctx.Args)
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)

	if argsLen != 1 {
		ctx.SendMsg("Usage: `!mod <default/spedup/slowed>`")
		return nil
	}
	mod := ctx.Args[0]
	if _, ok := Mods[mod]; !ok {
		ctx.SendMsg("Usage: `!mod <default/spedup/slowed>`")
		return nil
	}
	sess.SetMod(Mods[mod])
	return nil
}
