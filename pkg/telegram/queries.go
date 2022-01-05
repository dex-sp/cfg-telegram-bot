package telegram

import (
	"fmt"
	"io/ioutil"

	"github.com/dex-sp/cfg-telegram-bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	registrationQuery = "registration"
	locationQuery     = "location"
	priceQuery        = "price"
	payQuery          = "pay"
	orderQuery        = "order"
	callQuery         = "call"
	rulesQuery        = "rules"

	anotherDayQuery  = "another"
	changePhoneQuery = "change"
)

func (b *Bot) handleRegistrationQuery(query *tgbotapi.CallbackQuery) error {

	msgToPlayer := tgbotapi.NewMessage(query.From.ID, fmt.Sprintf(
		b.config.QueryResponses.Registration,
		query.From.FirstName, b.config.Owner.Name, b.config.Owner.CreditCard))
	msgToPlayer.ParseMode = "Markdown"

	playerPhone, err := b.userDataRepository.Get(query.From.ID, repository.Phones)
	if err != nil {
		return err
	}

	// First registration of Player
	if playerPhone == "" {
		msgToPlayer.Text = fmt.Sprintf(b.config.QueryResponses.FirstRegistration,
			getPhoneButton.Text, registrationButton.Text, payButton.Text, callButton.Text)

		msgToPlayer.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(getPhoneButton))

	} else {
		msgToPlayer.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(priceButton),
			tgbotapi.NewInlineKeyboardRow(changePhoneButton),
			tgbotapi.NewInlineKeyboardRow(locationButton))
	}

	_, err = b.bot.Send(msgToPlayer)
	return err
}

func (b *Bot) handleLocationQuery(query *tgbotapi.CallbackQuery) error {

	msg := tgbotapi.NewMessage(query.From.ID, fmt.Sprintf(
		b.config.QueryResponses.Location, b.config.LocationURL))
	msg.ParseMode = "Markdown"

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(registrationButton),
		tgbotapi.NewInlineKeyboardRow(rulesButton),
		tgbotapi.NewInlineKeyboardRow(mainChatButton))

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handlePriceQuery(query *tgbotapi.CallbackQuery) error {

	msg := tgbotapi.NewMessage(query.From.ID, b.config.QueryResponses.Price)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(registrationButton))

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handlePayQuery(query *tgbotapi.CallbackQuery) error {

	msg := tgbotapi.NewMessage(query.From.ID, b.config.QueryResponses.Price)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(registrationButton))

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleGetGameRulesQuery(query *tgbotapi.CallbackQuery) error {

	var err error
	fName := "rule_book.pdf"
	fileBytes, err := ioutil.ReadFile(fName)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewDocument(query.From.ID,
		tgbotapi.FileBytes{
			Name:  fName,
			Bytes: fileBytes})

	go func() {
		_, err = b.bot.Send(msg)
	}()
	return err
}
