package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func containsUserPhone(message *tgbotapi.Message) bool {
	return message.Contact.PhoneNumber != "" &&
		message.From.ID == message.Contact.UserID
}
