package cmd

import (
	"Gondon/bot"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	PlayerButtons = []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Emoji: discordgo.ComponentEmoji{
						Name: "⏹",
					},
					Style:    discordgo.SecondaryButton,
					CustomID: "button-stop",
				},
				discordgo.Button{
					Emoji: discordgo.ComponentEmoji{
						Name: "⏭",
					},
					Style:    discordgo.SecondaryButton,
					CustomID: "button-skip",
				},
				discordgo.Button{
					Emoji: discordgo.ComponentEmoji{
						Name: "🔂",
					},
					Style:    discordgo.SecondaryButton,
					CustomID: "button-repeat-one",
				},
				discordgo.Button{
					Emoji: discordgo.ComponentEmoji{
						Name: "💣",
					},
					Style:    discordgo.SecondaryButton,
					CustomID: "button-bassboost",
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "slowed&reverb",
					Style:    discordgo.SecondaryButton,
					CustomID: "button-slowed",
				},
				discordgo.Button{
					Label:    "default",
					Style:    discordgo.SecondaryButton,
					CustomID: "button-default",
				},
				discordgo.Button{
					Label:    "spedup!",
					Style:    discordgo.SecondaryButton,
					CustomID: "button-spedup",
				},
			},
		}}
	PlayerButtonsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.Context){
		"button-skip": func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.Context) {
			if !inVoiceChannel(ctx) {
				interactionRespondWithMessage(s, i, "Gandon is not in a voice channel!")
				return
			}
			err := SkipCommand(*ctx)
			if err != nil {
				ctx.SendMsg("Something went wrong...")
				log.WithFields(log.Fields{"error": err}).Warn()
			}
			InteractionRespondSilent(s, i)
		},
		"button-stop": func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.Context) {
			if !inVoiceChannel(ctx) {
				interactionRespondWithMessage(s, i, "Gandon is not in a voice channel!")
				return
			}
			err := StopCommand(*ctx)
			if err != nil {
				ctx.SendMsg("Something went wrong...")
				log.WithFields(log.Fields{"error": err}).Warn()
			}
			InteractionRespondSilent(s, i)
		},
		"button-bassboost": func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.Context) {
			if !inVoiceChannel(ctx) {
				interactionRespondWithMessage(s, i, "Gandon is not in a voice channel!")
				return
			}
			err := BassboostCommand(*ctx)
			if err != nil {
				ctx.SendMsg("Something went wrong...")
				log.WithFields(log.Fields{"error": err}).Warn()
			}
			InteractionRespondSilent(s, i)
		},
		"button-slowed": func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.Context) {
			if !inVoiceChannel(ctx) {
				interactionRespondWithMessage(s, i, "Gandon is not in a voice channel!")
				return
			}
			ctx.Args = []string{"slowed"}
			err := ModCommand(*ctx)
			if err != nil {
				ctx.SendMsg("Something went wrong...")
				log.WithFields(log.Fields{"error": err}).Warn()
			}
			interactionRespondWithMessage(s, i, "Set audio mod: slowed+reverb!")
		},
		"button-default": func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.Context) {
			if !inVoiceChannel(ctx) {
				interactionRespondWithMessage(s, i, "Gandon is not in a voice channel!")
				return
			}
			ctx.Args = []string{"default"}
			err := ModCommand(*ctx)
			if err != nil {
				ctx.SendMsg("Something went wrong...")
				log.WithFields(log.Fields{"error": err}).Warn()
			}
			interactionRespondWithMessage(s, i, "Set audio mod: default!")
		},
		"button-spedup": func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.Context) {
			if !inVoiceChannel(ctx) {
				interactionRespondWithMessage(s, i, "Gandon is not in a voice channel!")
				return
			}
			ctx.Args = []string{"spedup"}
			err := ModCommand(*ctx)
			if err != nil {
				ctx.SendMsg("Something went wrong...")
				log.WithFields(log.Fields{"error": err}).Warn()
			}
			interactionRespondWithMessage(s, i, "Set audio mod spedup!")
		},
		"button-repeat-one": func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.Context) {
			if !inVoiceChannel(ctx) {
				interactionRespondWithMessage(s, i, "Gandon is not in a voice channel!")
				return
			}
			sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
			if sess.RepeatMode == bot.RepeatModeNone {
				sess.RepeatMode = bot.RepeatModeOne
				interactionRespondWithMessage(s, i, "Repeat mode enabled!")
			} else {
				sess.RepeatMode = bot.RepeatModeNone
				interactionRespondWithMessage(s, i, "Repeat mode disabled")
			}
		},
	}
)

func InteractionRespondSilent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate, // no empty response in the API yet, using this type
	})
	if err != nil {
		panic(err)
	}
}
func interactionRespondWithMessage(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		panic(err)
	}
}
func interactionRespondWithInput(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Content: "TSTNG",
		},
	})
	if err != nil {
		panic(err)
	}
}
