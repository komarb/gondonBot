package bot

import "github.com/bwmarrin/discordgo"

const (
	RepeatModeNone int = 0
	RepeatModeOne      = 1
)

type Session struct {
	Queue         *SongQueue
	guildId       string
	ChannelId     string
	connection    *Connection
	PlayerMessage *discordgo.Message
	RepeatMode    int
}

type SessionManager struct {
	sessions map[string]*Session
}

type JoinProperties struct {
	Muted    bool
	Deafened bool
}
type PlayerMessage struct {
	MsgId     string
	ChannelId string
}

func NewSessionManager() *SessionManager {
	return &SessionManager{make(map[string]*Session)}
}

func newSession(guildId, channelId string, connection *Connection, playerMessage *discordgo.Message) *Session {
	session := new(Session)
	session.Queue = newSongQueue()
	session.guildId = guildId
	session.ChannelId = channelId
	session.connection = connection
	session.PlayerMessage = playerMessage
	session.RepeatMode = RepeatModeNone
	return session
}

func (sess Session) Play(song Song) error {
	return sess.connection.Play(song.Ffmpeg(sess.connection.mod))
}

func (sess Session) ToogleBassboost() bool {
	sess.connection.SetBassboost(!sess.connection.GetBassboost())
	return sess.connection.GetBassboost()
}

func (sess Session) SetMod(mod string) {
	sess.connection.SetMod(mod)
}

func (sess *Session) Stop() {
	sess.connection.Stop()
}

func (manager SessionManager) GetByGuild(guildId string) *Session {
	for _, sess := range manager.sessions {
		if sess.guildId == guildId {
			return sess
		}
	}
	return nil
}

func (manager SessionManager) GetByChannel(channelId string) (*Session, bool) {
	sess, found := manager.sessions[channelId]
	return sess, found
}

func (manager *SessionManager) Join(discord *discordgo.Session, guildId, channelId string,
	properties JoinProperties) (*Session, error) {
	vc, err := discord.ChannelVoiceJoin(guildId, channelId, properties.Muted, properties.Deafened)
	if err != nil {
		return nil, err
	}
	sess := newSession(guildId, channelId, NewConnection(vc), &discordgo.Message{})
	manager.sessions[channelId] = sess
	return sess, nil
}

func (manager *SessionManager) Leave(discord *discordgo.Session, session Session) {
	session.connection.Stop()
	session.connection.Disconnect()
	delete(manager.sessions, session.ChannelId)
}
