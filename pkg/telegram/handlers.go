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
	locationQuery     = "location"
	scheduleQuery     = "schedule"
	priceQuery        = "price"
	cmdPay            = "pay"
	orderQuery        = "order"
	callQuery         = "call"
	anotherDayQuery   = "another"
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

	case locationQuery:
		return b.handleLocationQuery(query)

	case scheduleQuery:
		return b.handleScheduleQuery(query)

	case priceQuery:
		return b.handlePriceQuery(query)

	case orderQuery:
		return b.handleOrderQuery(query)

	case callQuery:
		return b.handleCallQuery(query)

	case anotherDayQuery:

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

func (b *Bot) handleLocationQuery(query *tgbotapi.CallbackQuery) error {

	msg := tgbotapi.NewMessage(query.From.ID,
		"TODO: написать справку по локации")

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(registrationButton),
		tgbotapi.NewInlineKeyboardRow(chartButton))

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handlePriceQuery(query *tgbotapi.CallbackQuery) error {

	msg := tgbotapi.NewMessage(query.From.ID,
		"300р. - сыграть одну игру, примерно 40 минут\n"+
			"600р. - с 19:00 до 24:00\n"+
			"800р. - с 19:00 до 03:00\n\n"+
			"Среда")

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(registrationButton))

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleScheduleQuery(query *tgbotapi.CallbackQuery) error {

	msg := tgbotapi.NewMessage(query.From.ID,
		"Среда с 19:00 до 24:00\n"+
			"Пятница с 19:00 до 03:00\n"+
			"Суббота с 16:30 до 06:00")

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(registrationButton),
		tgbotapi.NewInlineKeyboardRow(anotherDayButton))

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleOrderQuery(query *tgbotapi.CallbackQuery) error {

	msg := tgbotapi.NewMessage(query.From.ID,
		"Напишите, пожалуйста, свой номер телефона. "+
			"С вами свяжется менеджер и вы обсудите условия.\n\n"+
			"Стоимость часа ведущего от 2800р.")

	if query.From.FirstName != "" {
		msg.Text = fmt.Sprintf("%s, напишите, пожалуйста, свой номер телефона. "+
			"С вами свяжется менеджер и вы обсудите условия.\n\n"+
			"Стоимость часа ведущего от 2800р.",
			query.From.FirstName)
	}

	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(getPhoneButton))

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleCallQuery(query *tgbotapi.CallbackQuery) error {

	msg := tgbotapi.NewMessage(query.From.ID,
		"Напишите свой номер телефона")

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(getPhoneButton),
		tgbotapi.NewKeyboardButtonRow(getLocationButton))

	msg.ReplyMarkup = keyboard

	_, err := b.bot.Send(msg)
	return err
}
