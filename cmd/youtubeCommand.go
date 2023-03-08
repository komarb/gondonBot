package cmd

import (
	"Gondon/bot"
	"bytes"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

const result_format = "\n`%d` $s - %s (%s)"

type ytSearchSession struct {
	results []bot.YTSearchContent
}
type ytSearchSessions map[string]ytSearchSession

var ytSessions ytSearchSessions = make(ytSearchSessions)

func ytSessionIdentifier(user *discordgo.User, channel *discordgo.Channel) string {
	return user.ID + channel.ID
}

func formatDuration(input string) string {
	return parseISO8601(input).String()
}

func YoutubeCommand(ctx bot.Context) {
	if len(ctx.Args) == 0 {
		ctx.SendMsg("usage `!yt <search query>`")
		return
	}

	results, err := YoutubeSearch(ctx)
	if err != nil {
		ctx.SendMsg(err.Error())
		return
	}
	if len(results) == 0 {
		ctx.SendMsg("No results found for your query `" + strings.Join(ctx.Args, " ") + "`.")
		return
	}

	buffer := bytes.NewBufferString("__Search results__ for `" + strings.Join(ctx.Args, " ") + "`:\n")
	for index, result := range results {
		buffer.WriteString(fmt.Sprintf(result_format, index+1, result.Snippet.Title, result.Snippet.ChannelTitle))
	}
	buffer.WriteString("\n\nTo pick a song, use `!pick <number>`.")
	ytSessions[ytSessionIdentifier(ctx.User, ctx.TextChannel)] = ytSearchSession{results}
	ctx.SendMsg(buffer.String())
}

func YoutubeSearch(ctx bot.Context) ([]bot.YTSearchContent, error) {
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if sess == nil {
		return nil, errors.New("Not in a voice channel!")
	}
	query := strings.Join(ctx.Args, " ")
	results, err := ctx.Youtube.Search(query)
	if err != nil {
		return nil, errors.New("Error searching YouTube: " + err.Error())
	}
	return results, nil
}
