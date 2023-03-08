package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type PlayerMessages struct {
	playerMessages map[string]*discordgo.Message
}

func NewPlayerMessages() *PlayerMessages {
	return &PlayerMessages{make(map[string]*discordgo.Message)}
}
func (pm *PlayerMessages) Update(guildID string, msg *discordgo.Message) {
	pm.playerMessages[guildID] = msg
}
func (pm PlayerMessages) Get(guildID string) *discordgo.Message {
	if _, ok := pm.playerMessages[guildID]; ok {
		return pm.playerMessages[guildID]
	}
	return nil
}

func NewPlayerEmbed(status string, track string, duration string, requestedBy string) *discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{}
	embed.Title = status
	embed.Color = 0x141414
	fields := make([]*discordgo.MessageEmbedField, 0)
	//fields[0].Name = "Track"
	//fields[0].Value = "TestTrack"
	//fields[0].Inline = true
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Track", Value: track, Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Duration", Value: fmt.Sprintf("`%s`", duration), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Requested by", Value: requestedBy, Inline: true})
	//fields[1].Name = "Requested by"
	//fields[1].Value = "124"
	//fields[1].Inline = true
	//
	//fields[2].Name = "Duration"
	//fields[2].Value = "`03:45`"
	//fields[2].Inline = true
	embed.Fields = append(embed.Fields, fields...)

	return &embed
}

func UpdatePlayerMsg(ctx *Context) error {
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	playerMessage := ctx.PlayerMessages.Get(ctx.Guild.ID)
	if playerMessage == nil {
		return fmt.Errorf("no Player message in this guild")
	}
	song := sess.Queue.Current()

	editedMsg := discordgo.NewMessageEdit(playerMessage.ChannelID, playerMessage.ID)
	if song == nil || sess.Queue.Running == false {
		editedMsg.SetEmbed(NewPlayerEmbed("🔇Stopped", "-", "-:-", "-"))
	} else {
		editedMsg.SetEmbed(NewPlayerEmbed("🔊Playing!", song.Title, song.Duration, "UpdatePlayerMsgTODO"))
	}
	_, err := ctx.Discord.ChannelMessageEditComplex(editedMsg)
	if err != nil {
		return fmt.Errorf("can't edit Player message: %s", err)
	}
	return nil
}
