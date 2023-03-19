package cmd

import (
	"Gondon/bot"
	"github.com/bwmarrin/discordgo"
)

func HelpCommand(ctx bot.Context) error {
	helpMsg := `- /help - Displays list of all available commands
	- /play [youtube URL|youtube search query] - Searches YouTube video/playlist by URL/search query and adds it to queue
	- /player - Makes bot send message with embedded music player
	- /pop - Remove last song from a queue
	- /mem - Generate meme (demotivator) with random messages and attachment from server. Reply to a message with attachment with this command to make meme with specified attachment.
	- /popusk - Finds out what user is a loser today`

	msgComplex := discordgo.MessageSend{Content: helpMsg}
	ctx.SendComplex(&msgComplex)
	return nil
}
