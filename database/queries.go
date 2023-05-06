package database

import (
	"chatbot/common"
	"chatbot/model"
	"context"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetBotAccount(collection *mongo.Collection) []model.Profile {
	match := bson.M{"meta.isBot": true}
	start := time.Now()
	common.Sugar.Infow("Getting the bots account",
		"database", collection.Database().Name(),
		"collection", collection.Name())
	cur, err := collection.Find(context.Background(), match)

	if err != nil {
		common.Sugar.Fatalw("Error getting the bot accounts",
			"database", collection.Database().Name(),
			"collection", collection.Name(),
			"duration", time.Since(start).String(),
			"error", err)
	}

	bots := []model.Profile{}

	err = cur.All(context.Background(), &bots)

	if err != nil {
		common.Sugar.Fatalw("Error getting the bot accounts",
			"database", collection.Database().Name(),
			"collection", collection.Name(),
			"duration", time.Since(start).String(),
			"error", err)
	}

	if len(bots) > 0 {
		common.Sugar.Infow("Successfully got the bot accounts",
			"database", collection.Database().Name(),
			"collection", collection.Name(),
			"duration", time.Since(start).String(),
			"count", len(bots))

		for i, bot := range bots {

			common.Sugar.Infow(
				"Bot "+strconv.Itoa(i+1)+": ",
				"ID", bot.ID.Hex(),
				"name", bot.Username,
				"gender", bot.Meta.Gender,
				"socketsCount", len(bot.Sockets))
		}
	} else {
		common.Sugar.Infow(
			"No bot accounts has been found",
			"database", collection.Database().Name(),
			"collection", collection.Name())
	}
	return bots
}

func GetConversation(collection *mongo.Collection, cid primitive.ObjectID) *model.Conversation {

	match := bson.M{
		"_id": cid,
	}

	var conversation model.Conversation
	start := time.Now()

	common.Sugar.Infow("Getting the conversation",
		"database", collection.Database().Name(),
		"collection", collection.Name(),
		"conversation", cid.Hex(),
		"duration", time.Since(start).String(),
	)

	err := collection.FindOne(context.Background(), match).Decode(&conversation)

	if err != nil {
		common.Sugar.Fatalw("Error getting conversation",
			"database", collection.Database().Name(),
			"collection", collection.Name(),
			"duration", time.Since(start).String(),
			"error", err)
	}

	common.Sugar.Infow("Conversation has been found",
		"database", collection.Database().Name(),
		"collection", collection.Name(),
		"conversation", cid.Hex(),
		"duration", time.Since(start).String(),
	)

	return &conversation
}

func ClearSocket(collection *mongo.Collection, pids []primitive.ObjectID) {

	match := bson.M{"_id": bson.M{"$nin": pids}}

	update := bson.M{"$set": bson.M{"sockets": []interface{}{}}}

	start := time.Now()

	_, err := collection.UpdateMany(context.TODO(), match, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			common.Sugar.Fatalw("Profiles not found",
				"database", collection.Database().Name(),
				"collection", collection.Name(),
				"duration", time.Since(start).String(),
				"error", err)
		}
		return
	}

}

func GetMessages(collection *mongo.Collection, cid primitive.ObjectID) *[]model.Message {

	match := bson.M{
		"conversation": cid,
	}

	start := time.Now()

	var message []model.Message

	findOptions := options.Find()

	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}})

	findOptions.SetLimit(10)

	common.Sugar.Infow("Getting the messages",
		"database", collection.Database().Name(),
		"collection", collection.Name(),
		"conversation", cid.Hex(),
		"duration", time.Since(start).String(),
	)

	cursor, err := collection.Find(context.Background(), match, findOptions)

	if err != nil {
		common.Sugar.Fatalw("Error getting messages",
			"database", collection.Database().Name(),
			"collection", collection.Name(),
			"duration", time.Since(start).String(),
			"error", err)
	}
	for cursor.Next(context.TODO()) {
		var msg model.Message
		_ = cursor.Decode(&msg)
		message = append(message, msg)
	}

	common.Sugar.Infow("Messages has been found",
		"database", collection.Database().Name(),
		"collection", collection.Name(),
		"conversation", cid.Hex(),
		"duration", time.Since(start).String(),
	)
	return &message
}

func GetProfileById(collection *mongo.Collection, pid primitive.ObjectID) (*model.Profile, error) {
	var profile model.Profile

	match := bson.M{"_id": pid}

	err := collection.FindOne(context.Background(), match).Decode(&profile)

	if err != nil {
		return nil, fmt.Errorf("NOT_FOUND")
	}

	return &profile, nil
}
