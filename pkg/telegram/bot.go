package telegram

import (
	"log"

	"github.com/dex-sp/cfg-telegram-bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot                *tgbotapi.BotAPI
	userDataRepository repository.UserDataRepository
}

func NewBot(bot *tgbotapi.BotAPI, repository repository.UserDataRepository) *Bot {

	return &Bot{bot: bot, userDataRepository: repository}
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

func containsUserPhone(message *tgbotapi.Message) bool {
	return message.Contact.PhoneNumber != "" &&
		message.From.ID == message.Contact.UserID
}

func (b *Bot) deleteReplyMenu(message *tgbotapi.Message) error {

	msg := tgbotapi.NewMessage(message.Chat.ID,
		"Спасибо, что помогаете нам улучшить качество сервиса.🥳")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

	_, err := b.bot.Send(msg)
	return err
}
