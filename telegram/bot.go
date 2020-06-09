package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	hopperApi "github.com/hopperteam/hopper-api/golang"
	"github.com/hopperteam/hopper-telegram-sp/hopper"
	"github.com/rs/zerolog/log"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	adapter *hopper.Adapter
}

func NewBot(token string, adapter *hopper.Adapter) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{bot, adapter}, nil
}

func (bot *Bot) EventLoop() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.bot.GetUpdatesChan(u)
	if err != nil {
		log.Err(err).Msg("Could not reply")
		return err
	}

	for update := range updates {
		if update.Message == nil || update.Message.Text == "" { // ignore any non-Message Updates
			continue
		}

		log.Info().
			Str("from", update.Message.From.UserName).
			Str("text", update.Message.Text).
			Int64("chatId", update.Message.Chat.ID).
			Msg("Got message")

		if update.Message.Command() == "subscribe" {
			go bot.handleSubscribePress(update.Message)
		} else {
			go bot.adapter.SendNotification(update.Message.Chat.ID, hopperApi.DefaultNotification(update.Message.From.FirstName + " " + update.Message.From.LastName, update.Message.Text))
		}
	}
	return nil
}

func (bot *Bot) handleSubscribePress(msg *tgbotapi.Message) {
	subR, err := bot.adapter.GetSubscribeLink(msg.Chat.ID, msg.Chat.Title)
	if err != nil {
		log.Err(err).Msg("Could not generate subscribe request")

		reply := tgbotapi.NewMessage(msg.Chat.ID, "Server Error: Could not generate subscribe request")
		_, err = bot.bot.Send(reply)
		if err != nil {
			log.Err(err).Msg("Could not reply")
			return
		}
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, "Subscribe to this chat")
	reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Subscribe", subR),
		),
	)

	_, err = bot.bot.Send(reply)
	if err != nil {
		log.Err(err).Msg("Could not reply")
		return
	}
}
