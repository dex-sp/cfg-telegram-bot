package telegram

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dex-sp/cfg-telegram-bot/pkg/config"
	"github.com/dex-sp/cfg-telegram-bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ledongthuc/pdf"
)

type Bot struct {
	bot                *tgbotapi.BotAPI
	userDataRepository repository.UserDataRepository
	config             *config.Config
}

var (
	// Start Menu buttons

	registrationButton tgbotapi.InlineKeyboardButton
	locationButton     tgbotapi.InlineKeyboardButton
	payButton          tgbotapi.InlineKeyboardButton
	priceButton        tgbotapi.InlineKeyboardButton
	callButton         tgbotapi.InlineKeyboardButton
	rulesButton        tgbotapi.InlineKeyboardButton
	mainChatButton     tgbotapi.InlineKeyboardButton

	// Phone & Location buttons

	getPhoneButton    tgbotapi.KeyboardButton
	getLocationButton tgbotapi.KeyboardButton

	// Payment button

	paymentConfirmedButton tgbotapi.InlineKeyboardButton
	paymentDeclinedButton  tgbotapi.InlineKeyboardButton

	//Other buttons

	changePhoneButton tgbotapi.InlineKeyboardButton

	startMenu tgbotapi.InlineKeyboardMarkup
)

func NewBot(bot *tgbotapi.BotAPI, userRepo repository.UserDataRepository, config *config.Config) *Bot {

	initButtons(config)
	return &Bot{bot: bot, userDataRepository: userRepo, config: config}
}

func (b *Bot) Start() error {

	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60

	updates := b.bot.GetUpdatesChan(upd)

	b.handleUpdates(updates)

	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {

	for update := range updates {
		if update.Message != nil { // If we got a message

			if update.Message.IsCommand() {
				if err := b.handleCommands(update.Message); err != nil {
					b.handleError(update.Message.Chat.ID, err)
				}
				continue
			}

			if err := b.handleMessage(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}

		} else if update.CallbackQuery != nil { // If we got a callback from button or event
			b.handleQueries(update.CallbackQuery)
		}
	}
}

func (b *Bot) deleteReplyMenu(message *tgbotapi.Message) error {

	msg := tgbotapi.NewMessage(message.Chat.ID, b.config.QueryResponses.Thanks)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) saveDocument(doc *tgbotapi.Document) (string, error) {

	docDirectURL, err := b.bot.GetFileDirectURL(doc.FileID)
	if err != nil {
		return "", err
	}

	//Get the response bytes from the url
	response, err := http.Get(docDirectURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", errors.New("received non 200 response code")
	}

	//Create a empty file
	filePath := doc.FileName
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	//Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

// Read PDF document as plain text.
func readDocument(path string, deleteAfterUse bool) (string, error) {

	doc, reader, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer doc.Close()

	resp, err := reader.GetPlainText()
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(resp)
	if err != nil {
		return "", err
	}

	if deleteAfterUse {
		err := doc.Close()
		if err != nil {
			return string(data), err
		}

		err = os.Remove(path)
		if err != nil {
			return string(data), err
		}
	}

	return string(data), nil
}

func (b *Bot) getUsersListString(rp repository.UserDataRepository, bk repository.Bucket) (string, error) {

	base := rp.GetAll(bk)
	var buffer bytes.Buffer
	var substring string
	var config tgbotapi.GetChatMemberConfig

	for userID, phone := range base {

		config.ChatConfigWithUser = tgbotapi.ChatConfigWithUser{
			ChatID: userID,
			UserID: userID,
		}
		user, err := b.bot.GetChatMember(config)
		if err != nil {
			return "", err
		}
		substring = fmt.Sprintf("User:\t@%s\nID:\t%d\nPhone:\t+%s\n\n",
			user.User.UserName, user.User.ID, phone)
		buffer.WriteString(substring)
	}
	return buffer.String(), nil
}

func initButtons(config *config.Config) {

	registrationButton = tgbotapi.NewInlineKeyboardButtonData(config.ButtonTemplates.Registration, registrationQuery)
	locationButton = tgbotapi.NewInlineKeyboardButtonData(config.ButtonTemplates.Location, locationQuery)
	priceButton = tgbotapi.NewInlineKeyboardButtonData(config.ButtonTemplates.Price, priceQuery)
	payButton = tgbotapi.NewInlineKeyboardButtonData(config.ButtonTemplates.Pay, payQuery)
	callButton = tgbotapi.NewInlineKeyboardButtonData(config.ButtonTemplates.Call, callQuery)
	rulesButton = tgbotapi.NewInlineKeyboardButtonData(config.ButtonTemplates.GameRules, rulesQuery)
	mainChatButton = tgbotapi.NewInlineKeyboardButtonURL(config.ButtonTemplates.MainChat, config.MainChat)

	paymentConfirmedButton = tgbotapi.NewInlineKeyboardButtonData(config.ButtonTemplates.PaymentConfirmed, confirmedPayment)
	paymentDeclinedButton = tgbotapi.NewInlineKeyboardButtonData(config.ButtonTemplates.PaymentDeclined, declinedPayment)

	changePhoneButton = tgbotapi.NewInlineKeyboardButtonData(config.ButtonTemplates.ChangePhone, changePhoneQuery)

	getPhoneButton = tgbotapi.NewKeyboardButtonContact(config.ButtonTemplates.GetPhone)
	getLocationButton = tgbotapi.NewKeyboardButtonLocation(config.ButtonTemplates.GetLocation)

	startMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(registrationButton),
		tgbotapi.NewInlineKeyboardRow(locationButton),
		tgbotapi.NewInlineKeyboardRow(priceButton),
		tgbotapi.NewInlineKeyboardRow(callButton),
		tgbotapi.NewInlineKeyboardRow(rulesButton),
		tgbotapi.NewInlineKeyboardRow(mainChatButton))
}
