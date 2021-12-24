package telegram

import (
	"fmt"

	"github.com/dex-sp/cfg-telegram-bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	registrationQuery = "registration"
	cancelQuery       = "cancel"
	locationQuery     = "location"
	priceQuery        = "price"
	cmdPay            = "pay"
	orderQuery        = "order"
	callQuery         = "call"

	anotherDayQuery  = "another"
	changePhoneQuery = "change"
)

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

	msg := tgbotapi.NewMessage(query.From.ID, b.config.QueryResponses.Cancel)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleLocationQuery(query *tgbotapi.CallbackQuery) error {

	msg := tgbotapi.NewMessage(query.From.ID,
		"TODO: написать справку по локации")

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(registrationButton),
		tgbotapi.NewInlineKeyboardRow(mainChartButton),
		tgbotapi.NewInlineKeyboardRow(guestChartButton))

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

func (b *Bot) handleCallQuery(query *tgbotapi.CallbackQuery) error {

	currentPhone, err := b.userDataRepository.Get(query.From.ID, repository.Phones)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(query.From.ID,
		fmt.Sprintf(b.config.QueryResponses.ChangePhone, changePhoneButton.Text))
	msg.ParseMode = "Markdown"

	if currentPhone == "" {
		msg.Text = b.config.QueryResponses.NewPhone
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(getPhoneButton))
	} else {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(changePhoneButton))
	}

	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleChangePhoneQuery(query *tgbotapi.CallbackQuery) error {

	msg := tgbotapi.NewMessage(query.From.ID, b.config.QueryResponses.SetPhone)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(getPhoneButton))

	_, err := b.bot.Send(msg)
	return err
}
