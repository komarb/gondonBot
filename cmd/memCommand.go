package cmd

import (
	"Gondon/bot"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"strings"
)

var (
	MemCommands = []*discordgo.ApplicationCommand{
		{
			Name: "мем",
			Type: discordgo.MessageApplicationCommand,
		},
	}
	MemCommandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.Context){
		"мем": func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.Context) {
			log.Info(ctx)
			err := MemCommand(*ctx)
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Warn("MemCommand error TODO: normal error")
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Something went wrong...",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					panic(err)
				}
			} else {
				InteractionRespondSilent(s, i)
			}
		},
	}
)

func MemCommand(ctx bot.Context) error {
	if ctx.Message.ReferencedMessage != nil {
		if len(ctx.Message.ReferencedMessage.Attachments) == 0 || !strings.Contains(ctx.Message.ReferencedMessage.Attachments[0].ContentType, "image") {
			if len(ctx.Message.ReferencedMessage.Embeds) == 0 || ctx.Message.ReferencedMessage.Embeds[0].Type != "gifv" {
				return nil
			}
		}
	}
	var url string
	var imgFormat string
	if ctx.Message.ReferencedMessage == nil {
		url, _ = getRandomImage(&ctx)
		imgFormat = url[strings.LastIndex(url, ".")+1:]
	} else if len(ctx.Message.ReferencedMessage.Attachments) != 0 {
		url = ctx.Message.ReferencedMessage.Attachments[0].URL
		imgFormat = url[strings.LastIndex(url, ".")+1:]
	} else if len(ctx.Message.ReferencedMessage.Embeds) != 0 {
		url = ctx.Message.ReferencedMessage.Embeds[0].URL
		imgFormat = "gif"
	} else {
		return nil
	}
	imgBytes, err := downloadFile(url)

	if err != nil {
		return fmt.Errorf("pic download failed")
	}

	bigText, err := getRandomText(&ctx, bot.BigText)
	if err != nil {
		return fmt.Errorf("retrieving big text failed: %s", err)
	}
	smallText, err := getRandomText(&ctx, bot.SmallText)
	if err != nil {
		return fmt.Errorf("retrieving small text failed: %s", err)
	}
	memImgBuf := bot.MakeDemotivator(imgBytes, imgFormat, bigText, smallText)
	if memImgBuf == nil {
		return fmt.Errorf("failed making demotivator for url: %s", url)
	}
	msg := discordgo.MessageSend{File: &discordgo.File{Name: "mem." + imgFormat, Reader: memImgBuf}}
	ctx.SendComplex(&msg)
	return nil
}

func getRandomText(ctx *bot.Context, variant int) (string, error) {
	matchStage := bson.D{{"$match", bson.D{{"$and",
		bson.A{
			bson.D{{"guildid", ctx.Guild.ID}},
			bson.D{{"texttype", variant}},
		},
	}}}}
	randomStage := bson.D{{"$sample", bson.D{{"size", 1}}}}
	cursor, err := ctx.TextsColl.Aggregate(context.TODO(), mongo.Pipeline{matchStage, randomStage})
	if err != nil {
		return "", fmt.Errorf("mongodb aggregation failed: %s", err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return "", fmt.Errorf("mongodb cursor decoding failed: %s", err)
	}
	return results[0]["content"].(string), nil
}

func getRandomImage(ctx *bot.Context) (string, error) {
	matchStage := bson.D{{"$match", bson.D{{"guildid", ctx.Guild.ID}}}}
	randomStage := bson.D{{"$sample", bson.D{{"size", 1}}}}
	cursor, err := ctx.ImgsColl.Aggregate(context.TODO(), mongo.Pipeline{matchStage, randomStage})
	if err != nil {
		return "", fmt.Errorf("mongodb aggregation failed: %s", err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return "", fmt.Errorf("mongodb cursor decoding failed: %s", err)
	}
	return results[0]["url"].(string), nil
}

func downloadFile(URL string) (*[]byte, error) {
	response, err := http.Get(URL)
	if err != nil {
		return nil, fmt.Errorf("can't get from url")
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("non 200 response")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read body")
	}
	return &body, nil
}

func SaveMessageToDb(msg *discordgo.Message, ctx *bot.Context) (int, error) {
	var err error
	var result *mongo.InsertManyResult
	// Images
	if len(msg.Attachments) != 0 {
		if !strings.Contains(msg.Attachments[0].ContentType, "image") {
			return 0, nil
		}
		imgUrl := msg.Attachments[0].URL
		_, err := ctx.ImgsColl.InsertOne(context.TODO(), bot.MemImg{URL: imgUrl, GuildID: ctx.Guild.ID})
		if err != nil {
			return 0, fmt.Errorf("error inserting image: %s", err)
		}
	} else { // Texts
		if len(msg.Content) == 0 || msg.Content[0] == '!' || msg.Author.ID != "333943790802960385" && msg.Author.ID != "259359326068670464" && msg.Author.ID != "684744749403340862" {
			return 0, nil
		}
		processedText := bot.ProcessText(msg.Content)
		if len(processedText) == 0 {
			return 0, nil
		}

		texts := make([]interface{}, 0)
		bigTextsStrings := bot.MakeSequences(processedText, bot.BigText)
		for _, str := range bigTextsStrings {
			texts = append(texts, bot.MemText{Content: str, AuthorID: msg.Author.ID, GuildID: ctx.Guild.ID, TextType: bot.BigText})
		}
		smallTextsStrings := bot.MakeSequences(processedText, bot.SmallText)
		for _, str := range smallTextsStrings {
			texts = append(texts, bot.MemText{Content: str, AuthorID: msg.Author.ID, GuildID: ctx.Guild.ID, TextType: bot.SmallText})
		}
		result, err = ctx.TextsColl.InsertMany(context.TODO(), texts)
		if err != nil {
			return 0, fmt.Errorf("failed to insert texts: %s", err)
		}
		return len(result.InsertedIDs), nil
	}
	return 0, nil
}
