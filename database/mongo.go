package database

import (
	"chatbot/common"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const user = "MongoUser"
const pass = "MongPass"
const url = "MongoUrl"
const dbname = "MongoDbName"

func Connect() *mongo.Database {

	uri := "mongodb+srv://" + user + ":" + pass + "@" + url + "/?authSource=admin"

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()
	common.Sugar.Infow("Try to connect to the mongo database", "dbname", dbname)
	start := time.Now()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		common.Sugar.Fatalw("Could not connect to the mongo database", "dbname", dbname, "error", err.Error(), "duration", time.Since(start).String())
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		common.Sugar.Fatalw("Could not connect to the mongo database", "dbname", dbname, "error", err.Error(), "duration", time.Since(start).String())
	}

	common.Sugar.Infow("Successfully connect to the database", "dbname", dbname, "duration", time.Since(start).String())

	db := client.Database(dbname)

	return db
}
