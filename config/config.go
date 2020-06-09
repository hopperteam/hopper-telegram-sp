package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

type Holder struct {
	DbConnectionStr string
	DbName string
	TelegramBotToken string
	BaseUrl string
}

var Config = LoadConfig()

func LoadConfig() *Holder {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return &Holder {
		DbConnectionStr: configStrReq("DB_CONNECTION_STRING"),
		DbName: configStrReq("DB_NAME"),
		TelegramBotToken: configStrReq("TELEGRAM_BOT_TOKEN"),
		BaseUrl: configStrReq("BASE_URL"),
	}
}

func configStrOpt(envName string, defaultValue string) string {
	val := os.Getenv(envName)
	if val == "" {
		return defaultValue
	}
	return val
}

func configStrReq(envName string) string {
	val := os.Getenv(envName)
	if val == "" {
		log.Fatal().Str("name", envName).Msg("Required env not specified")
	}
	return val
}
