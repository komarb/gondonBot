package bot

import (
	"bufio"
	"encoding/binary"
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/moutend/go-equalizer/pkg/equalizer"
	log "github.com/sirupsen/logrus"
	"io"
	"layeh.com/gopus"
	"os/exec"
)

const SIZE_OF_INT16 = 2
const CHANNELS int = 2
const FRAME_RATE int = 48000
const FRAME_SIZE int = 960
const MAX_BYTES = (FRAME_SIZE * CHANNELS * 2) * SIZE_OF_INT16

func (connection *Connection) sendPCM(voice *discordgo.VoiceConnection, pcm <-chan []int16) {
	connection.lock.Lock()
	if connection.sendpcm || pcm == nil {
		connection.lock.Unlock()
		return
	}
	connection.sendpcm = true
	connection.lock.Unlock()
	defer func() {
		connection.sendpcm = false
	}()
	encoder, err := gopus.NewEncoder(FRAME_RATE, CHANNELS, gopus.Audio)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Creating opus encoder error")
		return
	}
	for {
		receive, ok := <-pcm
		if !ok {
			log.Warning("PCM channel was closed")
			return
		}
		opus, err := encoder.Encode(receive, FRAME_SIZE, MAX_BYTES)
		if err != nil {
			log.Warning("Opus encoding error")
			return
		}
		if !voice.Ready || voice.OpusSend == nil {
			log.Warning("Lib discordgo is not ready for opus packets")
			return
		}
		voice.OpusSend <- opus
	}
}

func (connection *Connection) Play(ffmpeg *exec.Cmd) error {
	if connection.playing {
		return errors.New("song already playing")
	}
	connection.stopRunning = false
	out, err := ffmpeg.StdoutPipe()
	if err != nil {
		return err
	}
	buffer := bufio.NewReaderSize(out, 16384)
	err = ffmpeg.Start()
	if err != nil {
		return err
	}
	connection.playing = true
	defer func() {
		connection.playing = false
	}()
	connection.voiceConnection.Speaking(true)
	defer connection.voiceConnection.Speaking(false)
	if connection.send == nil {
		connection.send = make(chan []int16, 2)
	}
	go connection.sendPCM(connection.voiceConnection, connection.send)

	for {
		if connection.stopRunning {
			ffmpeg.Process.Kill()
			break
		}
		audioBuffer := make([]int16, FRAME_SIZE*CHANNELS)
		err = binary.Read(buffer, binary.LittleEndian, &audioBuffer)
		if connection.bassboost {
			audioBuffer = BassBoost(audioBuffer)
		}
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		if err != nil {
			return nil
		}
		connection.send <- audioBuffer
	}
	return nil
}

func (connection *Connection) Stop() {
	connection.stopRunning = true
	connection.playing = false
}

func BassBoost(data []int16) []int16 { //int16 -> float64 -> int16
	f := equalizer.NewPeaking(48000, 200, 50, 10)
	//f0 := equalizer.NewBandPass(48000, 440, 0.5) low quality filter
	for i := 0; i < len(data); i++ {
		var input float64
		input = Int16ToFloat(data[i])
		output := input
		output = f.Apply(output)
		data[i] = FloatToInt16(output)
	}
	return data
}

func FloatToInt16(f float64) int16 {
	var i int16
	f = f * 32768
	if f > 32767 {
		f = 32767
	}
	if f < -32768 {
		f = -32768
	}
	i = int16(f)
	return i
}

func Int16ToFloat(i int16) float64 {
	var f float64
	f = (float64(i)) / 32768.0
	if f > 1 {
		f = 1
	}
	if f < -1 {
		f = -1
	}
	return f
}
