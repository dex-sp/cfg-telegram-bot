package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (

	// Bot commands

	cmdStart = "start"

	// Bot queries

	registrationQuery = "registration"
	cancelQuery       = "cancel"
	cmdLocation       = "location"
	cmdSchedule       = "schedule"
	cmdPrice          = "price"
	cmdPay            = "pay"
	cmdOrder          = "order"
	callQuery         = "call"
)

func (b *Bot) handleCommands(message *tgbotapi.Message) error {

	switch message.Command() {

	case cmdStart:
		return b.handleStartCommand(message)

	default:
		return b.handleUnknownCommand(message)

	}
}

func (b *Bot) handleQueries(query *tgbotapi.CallbackQuery) error {

	switch query.Data {

	case registrationQuery:
		return b.handleRegistrationQuery(query)

	case cancelQuery:
		return b.handleCancelQuery(query)

	case callQuery:
		return b.handleCallQuery(query)

	default:

	}

	return nil
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {

	if containsUserPhone(message) {

		msg := tgbotapi.NewMessage(message.Chat.ID,
			fmt.Sprintf("[%s] %s", message.From.UserName, message.Contact.PhoneNumber))

		b.bot.Send(msg)

	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {

	msg := tgbotapi.NewMessage(message.Chat.ID,
		"Если у вас останутся вопросы - нажмите кнопу *Позвоните мне*, "+
			"и наш менеджер свяжемся с вами в ближайшее время!")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = startMenu

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {

	msg := tgbotapi.NewMessage(message.Chat.ID,
		"Извините, я не знаю такой команды 😔")
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleRegistrationQuery(query *tgbotapi.CallbackQuery) error {

	msg := tgbotapi.NewMessage(query.From.ID,
		"Напишите свой номер телефона")

	if query.From.FirstName != "" {
		msg.Text = fmt.Sprintf("%s, напишите свой номер телефона",
			query.From.FirstName)
	}

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleCancelQuery(query *tgbotapi.CallbackQuery) error {

	msg := tgbotapi.NewMessage(query.From.ID,
		"Ваша запись аннулирована. Будем рады вас видеть в другой день! 😉")

	if query.From.FirstName != "" {
		msg.Text = fmt.Sprintf("Окей, %s, ваша запись аннулирована. "+
			"Будем рады вас видеть в другой день! 😉",
			query.From.FirstName)
	}

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleCallQuery(query *tgbotapi.CallbackQuery) error {

	msg := tgbotapi.NewMessage(query.From.ID,
		"Напишите свой номер телефона")

	if query.From.FirstName != "" {
		msg.Text = fmt.Sprintf("%s, напишите свой номер телефона",
			query.From.FirstName)
	}

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(getPhoneButton),
		tgbotapi.NewKeyboardButtonRow(getLocationButton))

	msg.ReplyMarkup = keyboard

	_, err := b.bot.Send(msg)
	return err
}
