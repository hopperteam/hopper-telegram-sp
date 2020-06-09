package main

import (
	"github.com/hopperteam/hopper-telegram-sp/config"
	"github.com/hopperteam/hopper-telegram-sp/hopper"
	"github.com/hopperteam/hopper-telegram-sp/model"
	"github.com/hopperteam/hopper-telegram-sp/telegram"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting")

	log.Info().Msg("Connecting to Database")
	err := model.ConnectDB(config.Config.DbConnectionStr, config.Config.DbName)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to Database")
		return
	}
	log.Info().Msg("Connected to Database")

	log.Info().Msg("Connecting to Hopper")
	adapter, err := hopper.NewAdapter()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create hopperAdapter")
	}
	log.Info().Msg("Connected to Hopper")

	log.Info().Msg("Starting web server")
	go hopper.StartWebEndpoint()

	log.Info().Msg("Connecting to Telegram")
	bot, err := telegram.NewBot(config.Config.TelegramBotToken, adapter)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create bot")
		return
	}
	log.Info().Msg("Connected to Telegram")

	err = bot.EventLoop()
	log.Err(err).Msg("Error in event loop")
}
