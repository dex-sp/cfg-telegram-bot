package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (

	// Start Menu buttons

	registrationButton = tgbotapi.NewInlineKeyboardButtonData("📝 Записаться на игру", fmt.Sprintf("/%s", registrationQuery))
	cancelButton       = tgbotapi.NewInlineKeyboardButtonData("❌ Отменить запись", cancelQuery)
	locationButton     = tgbotapi.NewInlineKeyboardButtonData("📍 Место проведенеия игр", fmt.Sprintf("/%s", cmdLocation))
	scheduleButton     = tgbotapi.NewInlineKeyboardButtonData("🗓 Расписание", fmt.Sprintf("/%s", cmdSchedule))
	priceButton        = tgbotapi.NewInlineKeyboardButtonData("💵 Цены", fmt.Sprintf("/%s", cmdPrice))
	payButtom          = tgbotapi.NewInlineKeyboardButtonData("💳 Оплата", fmt.Sprintf("/%s", cmdPay))
	orderButton        = tgbotapi.NewInlineKeyboardButtonData("🎤 Заказать ведущего", fmt.Sprintf("/%s", cmdOrder))
	callButton         = tgbotapi.NewInlineKeyboardButtonData("📲 Позвоните мне", callQuery)
	chartButton        = tgbotapi.NewInlineKeyboardButtonURL("⚔️ В Чат", "http://1.com")

	// Phone & Location buttons

	getPhoneButton    = tgbotapi.NewKeyboardButtonContact("\xF0\x9F\x93\x9E Указать номер телефона")
	getLocationButton = tgbotapi.NewKeyboardButtonLocation("\xF0\x9F\x8C\x8F Указать геолокацию")

	// exitButton = tgbotapi.NewK

	startMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(registrationButton),
		tgbotapi.NewInlineKeyboardRow(cancelButton),
		tgbotapi.NewInlineKeyboardRow(locationButton),
		tgbotapi.NewInlineKeyboardRow(scheduleButton),
		tgbotapi.NewInlineKeyboardRow(priceButton),
		tgbotapi.NewInlineKeyboardRow(payButtom),
		tgbotapi.NewInlineKeyboardRow(orderButton),
		tgbotapi.NewInlineKeyboardRow(callButton),
		tgbotapi.NewInlineKeyboardRow(chartButton))
)
