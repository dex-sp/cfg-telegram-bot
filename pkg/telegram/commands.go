package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	startCommand = "start"
)

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {

	msg := tgbotapi.NewMessage(message.Chat.ID, b.config.CommandResponses.Start)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = startMenu

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {

	msg := tgbotapi.NewMessage(message.Chat.ID, b.config.CommandResponses.Default)
	_, err := b.bot.Send(msg)
	return err
}
