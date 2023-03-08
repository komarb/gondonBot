package cmd

import (
	"Gondon/bot"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func GetAllMsgCommand(ctx bot.Context) error {
	guild := ctx.Guild
	totalInserted := 0
	for _, ch := range guild.Channels {
		if ch.Type != 0 || ch.Name == "test-gondon" {
			continue
		}
		channelInserted := 0
		lastMsgID := ""
		log.Info(ch.Name)
		for {
			msgs, err := ctx.Discord.ChannelMessages(ch.ID, 100, lastMsgID, "", "")
			if err != nil || len(msgs) == 0 {
				break
			}
			for _, message := range msgs {
				insertedNum, err := SaveMessageToDb(message, &ctx)
				if err != nil {
					return fmt.Errorf("failed to insert: %s", err)
				}
				//if len(message.Content) == 0 || message.Content[0] == '!' || message.Author.ID != "333943790802960385" && message.Author.ID != "259359326068670464" && message.Author.ID != "684744749403340862" {
				//	continue
				//}
				//content := message.Content
				//content = bot.ProcessText(content)
				//if len(content) == 0 {
				//	continue
				//}
				//bigTextsStrings := bot.MakeSequences(content, bot.BigText)
				//texts := make([]interface{}, 0)
				//for _, str := range bigTextsStrings {
				//	texts = append(texts, bot.MemText{Content: str, AuthorID: message.Author.ID, GuildID: ctx.Guild.ID, TextType: bot.BigText})
				//}
				//smallTextsStrings := bot.MakeSequences(content, bot.SmallText)
				//for _, str := range smallTextsStrings {
				//	texts = append(texts, bot.MemText{Content: str, AuthorID: message.Author.ID, GuildID: ctx.Guild.ID, TextType: bot.SmallText})
				//}
				//result, err := ctx.TextsColl.InsertMany(context.TODO(), texts)
				//if err != nil {
				//	log.WithFields(log.Fields{"error": err}).Warning("Something went wrong while inserting...")
				//	return nil
				//}

				channelInserted += insertedNum
				totalInserted += insertedNum
			}
			log.Info("For channel ", ch.Name, " documents inserted: ", channelInserted)
			lastMsgID = msgs[len(msgs)-1].ID
		}
	}
	log.Info("Done! Total inserted: ", totalInserted)
	return nil
}
