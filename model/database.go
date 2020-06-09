package model

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var database *mongo.Database
var stateCollection *mongo.Collection
var chatSubscriptionsCollection *mongo.Collection

func ConnectDB(connectStr string, db string) error {
	dbOptions := options.Client().ApplyURI(connectStr)
	client, err := mongo.NewClient(dbOptions)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		return err
	}

	database = client.Database(db)
	stateCollection = database.Collection("state")
	chatSubscriptionsCollection = database.Collection("chatSubscriptions")

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_, err = stateCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{ "key": 1 },
	})

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_, err = chatSubscriptionsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{ "chat": 1 },
	})

	return err
}
