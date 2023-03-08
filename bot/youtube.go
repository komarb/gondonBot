package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

const (
	ERROR_TYPE    = -1
	VIDEO_TYPE    = 0
	PLAYLIST_TYPE = 1
)

type videoResponse struct {
	Formats []struct {
		Acodec     string  `json:"acodec"`
		SampleRate int     `json:"asr"`
		Url        string  `json:"url"`
		Quality    float64 `json:"quality"`
		FormatId   string  `json:"format_id"`
	} `json:"formats"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
}

type VideoResult struct {
	Media    string
	Title    string
	Duration string
}

type PlaylistVideo struct {
	Id       string `json:"id"`
	Playlist string `json:"playlist"`
}

type YTSearchContent struct {
	Id struct {
		VideoId string `json:"videoId"`
	} `json:"id"`
	Snippet struct {
		Title        string `json:"title"`
		Description  string `json:"description"`
		ChannelTitle string `json:"channelTitle"`
	} `json:"snippet"`
}

type ytApiResponse struct {
	Error   bool              `json:"error"`
	Content []YTSearchContent `json:"items"`
}

type Youtube struct {
	Cfg *Config
}

func (youtube Youtube) buildUrl(query string) (*string, error) {
	base := youtube.Cfg.ServiceUrl + "/youtube/v3/search"
	address, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("q", query)
	params.Add("part", "snippet")
	params.Add("type", "video")
	params.Add("key", youtube.Cfg.YtApiKey)
	address.RawQuery = params.Encode()
	str := address.String()
	return &str, nil
}
func (youtube Youtube) Search(query string) ([]YTSearchContent, error) {
	addr, err := youtube.buildUrl(query)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(*addr)
	if err != nil {
		return nil, err
	}
	var apiResp ytApiResponse
	json.NewDecoder(resp.Body).Decode(&apiResp)
	return apiResp.Content, nil
}
func (youtube Youtube) Get(input string) (int, *string, error) {
	cmdString := "./yt-dlp"
	if os.Getenv("GDN_DEBUG_MODE") == "1" {
		cmdString += "_macos"
	}
	cmd := exec.Command(cmdString, "--skip-download", "-j", "--no-simulate", input)
	//log.Info(cmdString, "--skip-download ", "-j ", "--no-simulate ", fmt.Sprintf("\"%s\"", input))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ERROR_TYPE, nil, err
	}
	str := out.String()
	//log.Info(str)
	return youtube.getType(str), &str, nil
}
func (youtube Youtube) getType(input string) int {
	if strings.Contains(input, "playlist_title") {
		return PLAYLIST_TYPE
	}
	if strings.Contains(input, "upload_date") {
		return VIDEO_TYPE
	}
	return ERROR_TYPE
}

func (youtube Youtube) Video(input string) (*VideoResult, error) {
	var resp videoResponse
	err := json.Unmarshal([]byte(input), &resp)
	if err != nil {
		return nil, err
	}
	durationString := fmt.Sprintf("%.2d:%.2d", resp.Duration/60, resp.Duration%60)
	for i := 0; i < len(resp.Formats); i++ {
		//if resp.Formats[i].Quality == 1.0 && resp.Formats[i].Acodec == "opus" && resp.Formats[i].SampleRate == 48000 {
		if resp.Formats[i].FormatId == "251" {
			return &VideoResult{resp.Formats[i].Url, resp.Title, durationString}, nil
		}
	}
	return nil, fmt.Errorf("needed codecs not found")
}

func (youtube Youtube) Playlist(input string) (*[]PlaylistVideo, string, error) {
	lines := strings.Split(input, "\n")
	videos := make([]PlaylistVideo, 0)
	playlistTitle := ""
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var video PlaylistVideo
		err := json.Unmarshal([]byte(line), &video)
		if err != nil {
			return nil, "", err
		}
		playlistTitle = video.Playlist
		videos = append(videos, video)
	}
	return &videos, playlistTitle, nil
}
