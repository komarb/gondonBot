package bot

import (
	"github.com/bwmarrin/discordgo"
	"sync"
)

type Connection struct {
	voiceConnection *discordgo.VoiceConnection
	send            chan []int16
	lock            sync.Mutex
	sendpcm         bool
	stopRunning     bool
	playing         bool
	bassboost       bool
	mod             string
}

func NewConnection(voiceConnection *discordgo.VoiceConnection) *Connection {
	connection := new(Connection)
	connection.voiceConnection = voiceConnection
	connection.bassboost = false
	connection.mod = "asetrate=48000"
	return connection
}

func (connection *Connection) Disconnect() {
	connection.voiceConnection.Disconnect()
}

func (connection *Connection) GetBassboost() bool {
	return connection.bassboost
}

func (connection *Connection) SetBassboost(value bool) {
	connection.bassboost = value
}

func (connection *Connection) GetMod() string {
	return connection.mod
}

func (connection *Connection) SetMod(mod string) {
	connection.mod = mod
}
