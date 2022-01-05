package telegram

import (
	"fmt"
	"strings"

	"github.com/dex-sp/cfg-telegram-bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func containsUserPhone(message *tgbotapi.Message) bool {
	return message.Contact != nil && message.Contact.PhoneNumber != "" &&
		message.From.ID == message.Contact.UserID
}

func (b *Bot) handleCallQuery(query *tgbotapi.CallbackQuery) error {

	playerPhone, err := b.userDataRepository.Get(query.From.ID, repository.Phones)
	if err != nil {
		return err
	}

	msgToPlayer := tgbotapi.NewMessage(query.From.ID,
		fmt.Sprintf(b.config.QueryResponses.ChangePhone, changePhoneButton.Text))
	msgToPlayer.ParseMode = "Markdown"

	// No phone number in our base
	if playerPhone == "" {
		msgToPlayer.Text = b.config.QueryResponses.NewPhone
		msgToPlayer.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(getPhoneButton))

	} else {
		msgToPlayer.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(changePhoneButton))

		// Message to owner about call request from player
		owmerNameSplit := strings.Split(b.config.Owner.Name, " ")
		msgToOwner := tgbotapi.NewMessage(b.config.Owner.TelegramID, fmt.Sprintf(
			b.config.QueryResponses.PlayerCallRequest,
			owmerNameSplit[0], query.From.UserName, playerPhone))

		_, err = b.bot.Send(msgToOwner)
		if err != nil {
			return err
		}
	}
	_, err = b.bot.Send(msgToPlayer)
	return err
}

func (b *Bot) handleChangePhoneQuery(query *tgbotapi.CallbackQuery) error {

	msg := tgbotapi.NewMessage(query.From.ID, b.config.QueryResponses.SetPhone)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(getPhoneButton))

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handlePhoneData(message *tgbotapi.Message) error {

	playerPhone, err := b.userDataRepository.Get(
		message.Contact.UserID, repository.Phones)
	if err != nil {
		return err
	}

	// Check phone number in our base
	if playerPhone != message.Contact.PhoneNumber {
		err := b.userDataRepository.Save(
			message.Contact.UserID,
			message.Contact.PhoneNumber,
			repository.Phones)
		if err != nil {
			return err
		}
	}

	// Delete reply "Get phone number" button
	if err := b.deleteReplyMenu(message); err != nil {
		return err
	}

	// Message to owner about new player
	if playerPhone == "" {
		owmerNameSplit := strings.Split(b.config.Owner.Name, " ")
		msg := tgbotapi.NewMessage(b.config.Owner.TelegramID, fmt.Sprintf(
			b.config.QueryResponses.FirstRegistrationOwnerNotification,
			owmerNameSplit[0], message.From.UserName,
			b.userDataRepository.Len(repository.Phones)))

		_, err = b.bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
