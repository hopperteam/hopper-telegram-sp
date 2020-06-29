package hopper

import (
	hopperApi "github.com/hopperteam/hopper-api/golang"
	"github.com/hopperteam/hopper-telegram-sp/config"
	"github.com/hopperteam/hopper-telegram-sp/model"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type Adapter struct {
	Api *hopperApi.HopperApi
	App *hopperApi.App
}

func NewAdapter() (*Adapter, error) {
	api := hopperApi.CreateHopperApi(hopperApi.HopperProd)
	ok, err := api.CheckConnectivity()
	if !ok {
		return nil, err
	}

	serialized, err := model.GetState("tgApp")
	if err == nil {
		app, err := api.DeserializeApp(serialized)
		if err == nil {
			log.Info().Str("id", app.Id).Msg("Deserialized App")
			return &Adapter{
				Api: api,
				App: app,
			}, nil
		}
	}

	log.Warn().Msg("Could not deserialize App, creating new App")
	app, err := api.CreateApp(
		"Telegram",
		config.Config.BaseUrl,
		config.Config.BaseUrl + "/logo.svg",
		config.Config.BaseUrl,
		"hello@hoppercloud.net",
	)
	if err != nil {
		return nil, err
	}

	serialized2, err := app.Serialize()
	if err != nil {
		return nil, err
	}

	err = model.SetState("tgApp", serialized2)
	if err != nil {
		return nil, err
	}

	log.Info().Str("id", app.Id).Msg("Created App")

	return &Adapter{
		Api: api,
		App: app,
	}, nil
}

func (adapter *Adapter) GetSubscribeLink(chatId int64, chatName string) (string, error) {
	return adapter.App.CreateSubscribeRequest(config.Config.BaseUrl + "/callback?chat=" + strconv.FormatInt(chatId, 10), &chatName)
}

func (adapter *Adapter) SendNotification(chatId int64, notification *hopperApi.Notification) {
	failed := make([]string, 0)
	log.Info().Int64("chatId", chatId).Interface("notification", notification).Msg("Sending Notification")
	err := model.GetSubscriptionsForChat(chatId, func(subId string) {
		_, err := adapter.Api.PostNotification(subId, notification)
		if err != nil {
			failed = append(failed, subId)
		}
	})

	if err != nil {
		log.Err(err).Msg("Could not get chat subscriptions")
		return
	}

	if len(failed) != 0 {
		adapter.retryFailed(failed, chatId, notification)
	}
}

func (adapter *Adapter) retryFailed(failed []string, chatId int64, notification *hopperApi.Notification) {
	time.Sleep(30 * time.Second)
	for _, subId := range failed {
		_, err := adapter.Api.PostNotification(subId, notification)
		if err != nil {
			log.Warn().Int64("chatId", chatId).Str("subscriptionId", subId).Err(err).Msg("Removing chat subscription after retry didn't work")
			err = model.DeleteChatSubscription(chatId, subId)
			if err != nil {
				log.Err(err).Int64("chatId", chatId).Str("subscriptionId", subId).Msg("Could not remove chat subscription")
			}
		}
	}
}
