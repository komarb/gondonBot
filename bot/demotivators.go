package bot

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

const (
	SmallText int = 0
	BigText       = 1
)

type BullyingToday struct {
	UserID string
	Curse  string
	Year   int
	Month  int
	Day    int
}

type MemText struct {
	Content  string
	AuthorID string
	GuildID  string
	TextType int
}
type MemImg struct {
	URL     string
	GuildID string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ProcessText(text string) string {
	re := regexp.MustCompile(`https?:\/\/.*[\s]*`)
	text = re.ReplaceAllString(text, "")
	excludeString := "[]{}:;'\"\\,<>/@#$%^&*_~=\n"
	for _, ch := range excludeString {
		text = strings.ReplaceAll(text, string(ch), "")
	}
	text = strings.ToLower(text)
	return text
}

func MakeSequences(text string, variant int) []string {
	res := make([]string, 0)
	ended := false
	words := strings.Split(text, " ")

	for i := 0; i < len(words); i++ {
		seqLength := 0
		if variant == BigText {
			seqLength = rand.Intn(2) + 1
		} else if variant == SmallText {
			seqLength = rand.Intn(7) + 1
		}
		seqString := ""
		for j := 0; j < seqLength; j++ {
			if i+j >= len(words) {
				ended = true
				break
			}
			seqString += words[i+j] + " "
		}
		if ended {
			break
		}
		for len(seqString) > 58 {
			lastSpace := strings.LastIndex(seqString, " ")
			if lastSpace == -1 {
				seqString = seqString[:58]
			} else {
				seqString = seqString[:lastSpace]
			}
		}
		if seqString != "" {
			seqString = strings.TrimRight(seqString, " ")
			if len(seqString) > 1 {
				res = append(res, seqString)
			}
		}
		i += seqLength + 1
	}
	if len(res) == 0 {
		res = append(res, text)
	}
	return res
}
