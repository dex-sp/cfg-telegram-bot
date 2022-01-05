package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommands(message *tgbotapi.Message) error {

	switch message.Command() {

	case startCommand:
		return b.handleStartCommand(message)

	case usersCommand:
		return b.handleUsersCommand(message)

	case playersCommand:
		return b.handlePlayersCommand(message)

	case game0verCommand:
		return b.handleGame0verCommand(message)

	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleQueries(query *tgbotapi.CallbackQuery) error {

	switch query.Data {

	case registrationQuery:
		return b.handleRegistrationQuery(query)

	case locationQuery:
		return b.handleLocationQuery(query)

	case priceQuery:
		return b.handlePriceQuery(query)

	case payQuery:
		return b.handlePayQuery(query)

	case callQuery:
		return b.handleCallQuery(query)

	case rulesQuery:
		return b.handleGetGameRulesQuery(query)

	case changePhoneQuery:
		return b.handleChangePhoneQuery(query)

	case confirmedPayment:
		return b.handleConfirmedPayment(query)

	case declinedPayment:
		return b.handleDeclinedPayment(query)
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
	} else {
		err := b.handleInvalidPaymentDocument(message)
		if err != nil {
			return err
		}
	}
	return nil
}
