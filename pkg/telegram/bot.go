package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {

	return &Bot{bot: bot}
}

func (b *Bot) Start() error {

	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60

	updates := b.bot.GetUpdatesChan(upd)

	b.handleUpdates(updates)

	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {

	for update := range updates {
		if update.Message != nil { // If we got a message

			if update.Message.IsCommand() {
				b.handleCommands(update.Message)
				continue
			}
			b.handleMessage(update.Message)

		} else if update.CallbackQuery != nil { // If we got a callback from button

			b.handleQueries(update.CallbackQuery)
		}
	}
}
