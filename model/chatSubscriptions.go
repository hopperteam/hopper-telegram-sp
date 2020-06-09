package model

import (
	"context"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ChatSubscription struct {
	Chat int64 `bson:"chat"`
	Subscription string `bson:"subscription"`
}

func InsertChatSubscription(chatId int64, subscription string) error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	_, err := chatSubscriptionsCollection.InsertOne(ctx, &ChatSubscription{
		Chat:         chatId,
		Subscription: subscription,
	})
	return err
}

func DeleteChatSubscription(chatId int64, subscription string) error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	_, err := chatSubscriptionsCollection.DeleteOne(ctx, &ChatSubscription{
		Chat:         chatId,
		Subscription: subscription,
	})
	return err
}

func GetSubscriptionsForChat(chatId int64, callback func(string)) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := chatSubscriptionsCollection.Find(ctx, bson.M{"chat": chatId})
	if err != nil {
		return err
	}

	chatSub := &ChatSubscription{}
	for cur.Next(ctx) {
		err := cur.Decode(chatSub)
		if err != nil {
			return err
		}
		callback(chatSub.Subscription)
	}

	return nil
}
