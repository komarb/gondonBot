package cmd

import (
	"Gondon/bot"
	"bufio"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"os"
	"time"
)

var popuskHistory = map[string]bot.BullyingToday{}
var cursesLoaded = false

func PopuskCommand(ctx bot.Context) error {
	users := getUsersInGuild(&ctx)
	log.Info(users)
	_, _, day := time.Now().Date()
	if !cursesLoaded {
		loadCurses(ctx.CursesColl)
		cursesLoaded = true
	}
	if popuskHistory[ctx.Guild.ID].Day != day {
		popusk := bot.BullyingToday{}
		rand.Seed(time.Now().UnixNano())
		popusk.UserID = users[rand.Intn(len(users))].ID
		popusk.Day = day
		curse, err := getRandomBullying(&ctx)
		if err != nil {
			return fmt.Errorf("retrieving curse word failed: %s", err)
		}
		popusk.Curse = curse
		popuskHistory[ctx.Guild.ID] = popusk
	}
	msg := discordgo.MessageSend{Content: "Попуск дня:\n" + "<@" + popuskHistory[ctx.Guild.ID].UserID + "> - ты **" + popuskHistory[ctx.Guild.ID].Curse + "**"}
	ctx.SendComplex(&msg)
	return nil
}

func getUsersInGuild(ctx *bot.Context) []discordgo.User {
	members, _ := ctx.Discord.GuildMembers(ctx.Guild.ID, "", 1000)
	res := make([]discordgo.User, 0)
	for _, member := range members {
		if !member.User.Bot {
			res = append(res, *member.User)
		}
	}
	return res
}

func loadCurses(coll *mongo.Collection) {
	data, err := os.Open("res/curses.txt")
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("Can't read curses file, shutting down...")
	}
	defer data.Close()

	fileScanner := bufio.NewScanner(data)
	fileScanner.Split(bufio.ScanLines)

	opts := options.Update().SetUpsert(true)
	var filter bson.D
	var update bson.D

	for fileScanner.Scan() {
		curse := fileScanner.Text()
		filter = bson.D{{"curse", curse}}
		update = bson.D{{"$set", bson.D{{"curse", curse}}}}
		_, err = coll.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Warn("Something went wrong while upserting")
		}
	}
}

func getRandomBullying(ctx *bot.Context) (string, error) {
	randomStage := bson.D{{"$sample", bson.D{{"size", 1}}}}
	cursor, err := ctx.CursesColl.Aggregate(context.TODO(), mongo.Pipeline{randomStage})
	if err != nil {
		return "", fmt.Errorf("mongodb aggregation failed: %s", err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return "", fmt.Errorf("mongodb cursor decoding failed: %s", err)
	}
	return results[0]["curse"].(string), nil
}
