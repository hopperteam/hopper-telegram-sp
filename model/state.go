package model

import (
	"context"
	hopperApi "github.com/hopperteam/hopper-api/golang"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type StateElement struct {
	Key string `bson:"key"`
	Value string `bson:"value"`
}

func GetState(keyName string) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	el := &StateElement{}
	err := stateCollection.FindOne(ctx, bson.M{"key": keyName}).Decode(el)
	if err != nil {
		return "", nil
	}

	return el.Value, nil
}

func SetState(keyName string, value string) error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	el := &StateElement{
		Key: keyName,
		Value: value,
	}

	_, err := stateCollection.ReplaceOne(ctx, bson.M{"key": keyName}, el, &options.ReplaceOptions{
		Upsert: hopperApi.BoolPtr(true),
	})
	return err
}
