package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (

	// Bot commands

	cmdStart = "start"

	// Bot queries

	cmdRegistration = "registration"
	cmdCancel       = "cancel"
	cmdLocation     = "location"
	cmdSchedule     = "schedule"
	cmdPrice        = "price"
	cmdPay          = "pay"
	cmdOrder        = "order"
	cmdCall         = "call"
)

func (b *Bot) handleCommands(message *tgbotapi.Message) error {

	switch message.Command() {
	case cmdStart:
		err := b.handleStartCommand(message)
		if err != nil {
			return err
		}
	default:
		err := b.handleUnknownCommand(message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) handleQueries(query *tgbotapi.CallbackQuery) error {

	switch query.Data {
	case cmdRegistration:
		err := b.handleCancelQuery(query)
		if err != nil {
			return err
		}
	case cmdCancel:
		err := b.handleCancelQuery(query)
		if err != nil {
			return err
		}
	default:

	}

	return nil
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {

	log.Printf("[%s] %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)

	b.bot.Send(msg)
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
