package telegram

import (
	"errors"
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
	cancelButton       tgbotapi.InlineKeyboardButton
	locationButton     tgbotapi.InlineKeyboardButton
	priceButton        tgbotapi.InlineKeyboardButton
	callButton         tgbotapi.InlineKeyboardButton
	mainChartButton    tgbotapi.InlineKeyboardButton
	guestChartButton   tgbotapi.InlineKeyboardButton

	// Phone & Location buttons

	getPhoneButton    tgbotapi.KeyboardButton
	getLocationButton tgbotapi.KeyboardButton

	//Other buttons

	changePhoneButton tgbotapi.InlineKeyboardButton

	startMenu tgbotapi.InlineKeyboardMarkup
)

func NewBot(bot *tgbotapi.BotAPI, repository repository.UserDataRepository, config *config.Config) *Bot {

	registrationButton = tgbotapi.NewInlineKeyboardButtonData(config.ButtonTemplates.Registration, registrationQuery)
	cancelButton = tgbotapi.NewInlineKeyboardButtonData(config.ButtonTemplates.Cancel, cancelQuery)
	locationButton = tgbotapi.NewInlineKeyboardButtonData(config.ButtonTemplates.Location, locationQuery)
	priceButton = tgbotapi.NewInlineKeyboardButtonData(config.ButtonTemplates.Price, priceQuery)
	callButton = tgbotapi.NewInlineKeyboardButtonData(config.ButtonTemplates.Call, callQuery)
	mainChartButton = tgbotapi.NewInlineKeyboardButtonURL(config.ButtonTemplates.MainChat, config.MainChat)
	guestChartButton = tgbotapi.NewInlineKeyboardButtonURL(config.ButtonTemplates.GuestChat, config.GuestChat)

	changePhoneButton = tgbotapi.NewInlineKeyboardButtonData(config.ButtonTemplates.ChangePhone, changePhoneQuery)

	getPhoneButton = tgbotapi.NewKeyboardButtonContact(config.ButtonTemplates.GetPhone)
	getLocationButton = tgbotapi.NewKeyboardButtonLocation(config.ButtonTemplates.GetLocation)

	startMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(registrationButton),
		tgbotapi.NewInlineKeyboardRow(cancelButton),
		tgbotapi.NewInlineKeyboardRow(locationButton),
		tgbotapi.NewInlineKeyboardRow(priceButton),
		tgbotapi.NewInlineKeyboardRow(callButton),
		tgbotapi.NewInlineKeyboardRow(mainChartButton),
		tgbotapi.NewInlineKeyboardRow(guestChartButton))

	return &Bot{bot: bot, userDataRepository: repository, config: config}
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

func containsUserPhone(message *tgbotapi.Message) bool {
	return message.Contact != nil && message.Contact.PhoneNumber != "" &&
		message.From.ID == message.Contact.UserID
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

	//Write the bytes to the fiel
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
		err := os.Remove(path)
		if err != nil {
			return string(data), err
		}
	}

	return string(data), nil
}
