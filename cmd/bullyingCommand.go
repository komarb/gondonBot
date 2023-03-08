package cmd

import (
	"Gondon/bot"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"time"
)

func PopuskCommand(ctx bot.Context) error {
	users := [...]string{"333943790802960385"}
	_, _, day := time.Now().Date()
	log.Info(day)
	log.Info(ctx.BullyingToday)
	if ctx.BullyingToday.Day != day {
		rand.Seed(time.Now().UnixNano())
		ctx.BullyingToday.UserID = users[rand.Intn(2)]
		ctx.BullyingToday.Day = day
		curse, err := getRandomBullying(&ctx)
		if err != nil {
			return fmt.Errorf("retrieving curse word failed: %s", err)
		}
		ctx.BullyingToday.Curse = curse
	}
	ctx.SendMsg("Попуск дня:\n" + "<@" + ctx.BullyingToday.UserID + "> - ты **" + ctx.BullyingToday.Curse + "**")
	return nil
}

func getRandomBullying(ctx *bot.Context) (string, error) {
	//matchStage := bson.D{{"$match", bson.D{{"guildid", ctx.Guild.ID}}}}
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
