package telegram

import (
	"strings"

	"github.com/dex-sp/cfg-telegram-bot/pkg/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func isPayment(data string, owner config.Owner) bool {

	cardSplit := strings.Split(owner.CreditCard, " ")
	lastDigits := cardSplit[len(cardSplit)-1] // 4 last digits of owner credit card

	nameSplit := strings.Split(owner.Name, " ")
	name := nameSplit[0]
	patronymic := nameSplit[1]
	surname := nameSplit[2]

	return strings.Contains(data, owner.Name) ||
		strings.Contains(data, owner.CreditCard) ||
		strings.Contains(data, owner.Phone) ||
		strings.Contains(data, lastDigits) ||
		strings.Contains(data, surname) ||
		(strings.Contains(data, name) && strings.Contains(data, patronymic))
}

func (b *Bot) handlePayment(message *tgbotapi.Message) error {

	return nil
}
