package cmd

import (
	"Gondon/bot"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)
import "fmt"
import "strings"

func PlayCommand(ctx bot.Context) error {
	videoURLS := make([]string, 0)
	if err := JoinCommand(ctx); err != nil {
		return nil
	}
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if len(ctx.Args) == 0 {
		ctx.SendMsg("Play command usage: `!play <youtube_url/youtube video name>'`")
		return nil
	}
	if strings.Contains(ctx.Args[0], "youtube.com") {
		videoURLS = append(videoURLS, ctx.Args...)
	} else {
		results, err := YoutubeSearch(ctx)
		if err != nil {
			return fmt.Errorf("error searching youtube: %s", err)
		}
		if len(results) == 0 {
			ctx.SendMsg("No results found for your query `" + strings.Join(ctx.Args, " ") + "`.")
			return nil
		}
		videoURLS = append(videoURLS, "https://www.youtube.com/watch?v="+results[0].Id.VideoId)
	}

	msg := ctx.SendMsg("Adding songs to queue...")
	for _, arg := range videoURLS {
		t, inp, err := ctx.Youtube.Get(arg)
		if err != nil {
			return fmt.Errorf("can't parse youtube request: %s", err)
		}

		switch t {
		case bot.ERROR_TYPE:
			return fmt.Errorf("got youtube-dl error type: %d", t)
		case bot.VIDEO_TYPE:
			video, err := ctx.Youtube.Video(*inp)
			if err != nil {
				return fmt.Errorf("can't get youtube video: %s", err)
			}
			song := bot.NewSong(video.Media, video.Title, video.Duration, arg)
			sess.Queue.Push(*song)
			ctx.Discord.ChannelMessageEdit(ctx.TextChannel.ID, msg.ID, "Added `"+song.Title+"` to the song queue.")
			break
		case bot.PLAYLIST_TYPE:
			videos, playlistTitle, err := ctx.Youtube.Playlist(*inp)
			if err != nil {
				return fmt.Errorf("can't get youtube playlist: %s", err)
			}
			for _, v := range *videos {
				id := v.Id
				_, i, err := ctx.Youtube.Get(id)
				if err != nil {
					log.WithFields(log.Fields{"error": err}).Warn("Can't get video id")
					continue
				}
				video, err := ctx.Youtube.Video(*i)
				if err != nil {
					return fmt.Errorf("can't get youtube video in playlist: %s", err)
				}
				song := bot.NewSong(video.Media, video.Title, video.Duration, arg)
				sess.Queue.Push(*song)
			}
			ctx.Discord.ChannelMessageEdit(ctx.TextChannel.ID, msg.ID, "Added playlist `"+playlistTitle+"` to the song queue.")
			break
		}
	}
	queue := sess.Queue
	if !queue.Running {
		go queue.Start(sess, func(track string, duration string, requestedBy string) {
			if pm := ctx.PlayerMessages.Get(ctx.Guild.ID); pm != nil {
				ctx.Discord.ChannelMessageDelete(pm.ChannelID, pm.ID)
			}
			embed := bot.NewPlayerEmbed("🔊Playing!", track, duration, ctx.User.Mention())
			msgComplex := discordgo.MessageSend{Embed: embed, Components: PlayerButtons}
			newPlayerMessage := ctx.SendComplex(&msgComplex)
			ctx.PlayerMessages.Update(ctx.Guild.ID, newPlayerMessage)
		})
	}
	return nil
}
