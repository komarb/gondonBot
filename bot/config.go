package bot

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	Prefix       string `json:"prefix"`
	ServiceUrl   string `json:"service_url"`
	BotToken     string `json:"bot_token"`
	GameStatus   string `json:"game_status"`
	YtApiKey     string `json:"yt_api_key"`
	DBUri        string `json:"db_uri"`
	FFmpegStatic string `json:"ffmpeg_static"`
}

func GetConfigEnv() *Config {
	var config Config
	config.Prefix = os.Getenv("GDN_PREFIX")
	config.ServiceUrl = os.Getenv("GDN_SERVICE_URL")
	config.BotToken = os.Getenv("GDN_BOT_TOKEN")
	config.GameStatus = os.Getenv("GDN_GAME_STATUS")
	config.YtApiKey = os.Getenv("GDN_YOUTUBE_API_KEY")
	config.DBUri = os.Getenv("GDN_DB_URI")
	config.FFmpegStatic = os.Getenv("GDN_FFMPEG_STATIC")
	return &config
}

func GetConfigFile() *Config {
	var config Config
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Can't read config.json file, shutting down...")
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Can't correctly parse json from config.json, shutting down...")
	}
	return &config
}
