package telegram

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dex-sp/cfg-telegram-bot/pkg/config"
	"github.com/dex-sp/cfg-telegram-bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	validPayment     string = "valid"
	invalidPayment   string = "invalid"
	confirmedPayment string = "confirmed"
	declinedPayment  string = "declined"
)

type Payment struct {
	From      *tgbotapi.User
	Document  *tgbotapi.Document
	Phone     string
	Confirmed bool
}

func isPayment(data string, owner config.Owner) bool {

	cardSplit := strings.Split(owner.CreditCard, " ")
	nameSplit := strings.Split(owner.Name, " ")

	return strings.Contains(data, owner.Name) ||
		strings.Contains(data, owner.CreditCard) ||
		strings.Contains(data, owner.Phone) ||
		strings.Contains(data, cardSplit[len(cardSplit)-1]) || // 4 last digits of owner credit card
		strings.Contains(data, nameSplit[2]) || // surname
		(strings.Contains(data, nameSplit[0]) &&
			strings.Contains(data, nameSplit[1])) // name & patronymic
}

/* Если предварительная проверка квитанции об оплате пройдена успешно, то
хозяину бота отправляется запрос на проверку банковского счета*/
func (b *Bot) handlePayment(message *tgbotapi.Message) error {

	confirmation, err := b.userDataRepository.Get(message.From.ID, repository.Confirmations)
	if err != nil {
		return err
	}

	// If the player is a confirmed participant in the nearest game
	if confirmation != "" {
		msgToPlayer := tgbotapi.NewMessage(message.From.ID, b.config.Errors.AlreadyConfirmed)
		_, err := b.bot.Send(msgToPlayer)
		if err != nil {
			return err
		}

	} else {
		// File Direct URL can be seen only by the owner of the bot
		fileDirectURL, err := b.bot.GetFileDirectURL(message.Document.FileID)
		if err != nil {
			return err
		}

		//Message to owner
		owmerNameSplit := strings.Split(b.config.Owner.Name, " ")
		msgToOwner := tgbotapi.NewMessage(b.config.Owner.TelegramID,
			fmt.Sprintf(b.config.QueryResponses.CheckPayment,
				owmerNameSplit[0], message.From.UserName, message.From.ID))

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(paymentConfirmedButton),
			tgbotapi.NewInlineKeyboardRow(paymentDeclinedButton),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(
					b.config.ButtonTemplates.GetPaymentDoc, fileDirectURL)))

		msgToOwner.ReplyMarkup = keyboard

		_, err = b.bot.Send(msgToOwner)
		if err != nil {
			return err
		}
	}
	return nil
}

/*Если вместо квитанции пришло нечто другое, то пользователю, который отправил эту
квитанцию прийдет соответсвтующее сообщение*/
func (b *Bot) handleInvalidPaymentDocument(message *tgbotapi.Message) error {

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(
		b.config.Errors.InvalidPaymentDocument, callButton.Text))
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(callButton))

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleConfirmedPayment(query *tgbotapi.CallbackQuery) error {

	playerID, err := userIDFromText(query.Message.Text)
	if err != nil {
		return err
	}

	msgToPlayer := tgbotapi.NewMessage(playerID, b.config.QueryResponses.PlayerNotification)
	msgToPlayer.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(locationButton))

	playerPhone, err := b.userDataRepository.Get(playerID, repository.Phones)
	if err != nil {
		return err
	}

	var msgToOwnerText string
	if playerPhone != "" {
		msgToOwnerText = fmt.Sprintf(b.config.QueryResponses.OwnerСonfirmedPayment,
			query.From.UserName, playerPhone)
	} else {
		msgToOwnerText = fmt.Sprintf(b.config.QueryResponses.OwnerСonfirmedPayment,
			query.From.UserName, "<не указан>")
	}

	msgToOwner := tgbotapi.NewMessage(b.config.Owner.TelegramID, msgToOwnerText)

	err = b.userDataRepository.Save(
		query.From.ID,
		time.Now().Format("2006-01-02 15:04:05"),
		repository.Confirmations)
	if err != nil {
		return err
	}

	if _, err := b.bot.Send(msgToPlayer); err != nil {
		return err
	}
	if _, err := b.bot.Send(msgToOwner); err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleDeclinedPayment(query *tgbotapi.CallbackQuery) error {

	playerID, err := userIDFromText(query.Message.Text)
	if err != nil {
		return err
	}

	msgToPlayer := tgbotapi.NewMessage(playerID, b.config.Errors.DeclinedPayment)
	msgToPlayer.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(callButton))

	playerPhone, err := b.userDataRepository.Get(playerID, repository.Phones)
	if err != nil {
		return err
	}

	var msgToOwnerText string
	if playerPhone != "" {
		msgToOwnerText = fmt.Sprintf(b.config.QueryResponses.OwnerDeclinedPayment,
			query.From.UserName, playerPhone)
	} else {
		msgToOwnerText = fmt.Sprintf(b.config.QueryResponses.OwnerDeclinedPayment,
			query.From.UserName, "<не указан>")
	}

	msgToOwner := tgbotapi.NewMessage(b.config.Owner.TelegramID, msgToOwnerText)

	if _, err := b.bot.Send(msgToPlayer); err != nil {
		return err
	}
	if _, err := b.bot.Send(msgToOwner); err != nil {
		return err
	}
	return nil
}

func userIDFromText(text string) (int64, error) {

	if !strings.Contains(text, "ID_") {
		return -1, errors.New("user ID not found in message")
	}

	re, err := regexp.Compile(`ID_(\d+)`)
	if err != nil {
		return -1, err
	}

	values := re.FindStringSubmatch(text)
	valLen := len(values)
	if valLen < 2 {
		return -1, errors.New("invalid user ID")
	}

	userID, err := strconv.ParseInt(values[valLen-1], 10, 64)
	if err != nil {
		return -1, err
	}
	return userID, nil
}
