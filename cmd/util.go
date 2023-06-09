package cmd

import (
	"Gondon/bot"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func parseISO8601(str string) time.Duration {
	durationRegex := regexp.MustCompile(`P(?P<years>\d+Y)?(?P<months>\d+M)?(?P<days>\d+D)?T?(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+S)?`)
	matches := durationRegex.FindStringSubmatch(str)
	years := parseInt64(matches[1])
	months := parseInt64(matches[2])
	days := parseInt64(matches[3])
	hours := parseInt64(matches[4])
	minutes := parseInt64(matches[5])
	seconds := parseInt64(matches[6])
	hour := int64(time.Hour)
	minute := int64(time.Minute)
	second := int64(time.Second)
	return time.Duration(years*24*365*hour + months*30*24*hour + days*24*hour + hours*hour + minutes*minute + seconds*second)
}

func parseInt64(value string) int64 {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(value[:len(value)-1])
	if err != nil {
		return 0
	}
	return int64(parsed)
}

func inVoiceChannel(ctx *bot.Context) bool {
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if sess == nil {
		return false
	}
	return true
}

func GetDbName(dburi string) string {
	l := 0
	i := 0
	for i != 3 {
		l += strings.Index(dburi[l:], "/") + 1
		i++
	}
	r := strings.Index(dburi, "?")
	if r == -1 {
		r = len(dburi)
	}
	return dburi[l:r]
}
