package bot

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strconv"
)

type Song struct {
	Media    string
	Title    string
	Duration string
	Id       string
}

func NewSong(media, title, duration, id string) *Song {
	song := new(Song)
	song.Media = media
	song.Title = title
	song.Duration = duration
	song.Id = id
	return song
}

func (song Song) Ffmpeg(mod string) *exec.Cmd {
	cmdArgs := []string{"-reconnect", "1", "-reconnect_at_eof", "1", "-reconnect_streamed", "1", "-reconnect_delay_max", "2", "-i", fmt.Sprintf("%s", song.Media), "-f", "s16le", "-af", mod, "-ac", strconv.Itoa(CHANNELS), "pipe:1"}
	cmdString := "ffmpeg"
	if os.Getenv("GDN_FFMPEG_STATIC") == "1" {
		cmdString = "./" + cmdString
	}
	if os.Getenv("GDN_DEBUG_MODE") == "1" {
		cmdString = cmdString + "_macos"
	}
	log.Info(cmdString, cmdArgs)
	return exec.Command(cmdString, cmdArgs...)
}

//"-filter:a asetrate=r=66000"
//"-af", "asetrate=48000,aresample=48000"
