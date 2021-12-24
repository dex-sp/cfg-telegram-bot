package telegram

import (
	"github.com/dex-sp/cfg-telegram-bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommands(message *tgbotapi.Message) error {

	switch message.Command() {

	case startCommand:
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

	case priceQuery:
		return b.handlePriceQuery(query)

	case callQuery:
		return b.handleCallQuery(query)

	case anotherDayQuery:

	case changePhoneQuery:
		return b.handleChangePhoneQuery(query)

	default:

	}

	return nil
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {

	if containsUserPhone(message) {
		return b.handlePhoneData(message)

	} else if message.Document != nil {
		return b.handleDocument(message)
	}

	return nil
}

func (b *Bot) handlePhoneData(message *tgbotapi.Message) error {

	currentPhone, err := b.userDataRepository.Get(message.Contact.UserID, repository.Phones)
	if err != nil {
		return err
	}

	if currentPhone != message.Contact.PhoneNumber {
		err := b.userDataRepository.Save(
			message.Contact.UserID,
			message.Contact.PhoneNumber,
			repository.Phones)
		if err != nil {
			return err
		}
	}
	return b.deleteReplyMenu(message)
}

func (b *Bot) handleDocument(message *tgbotapi.Message) error {

	docPath, err := b.saveDocument(message.Document)
	if err != nil {
		return err
	}

	data, err := readDocument(docPath, true)
	if err != nil {
		return err
	}

	if isPayment(data, b.config.Owner) {
		err := b.handlePayment(message)
		if err != nil {
			return err
		}
	}

	return nil
}
