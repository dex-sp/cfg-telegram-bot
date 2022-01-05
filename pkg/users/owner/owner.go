package owner

import (
	"github.com/dex-sp/cfg-telegram-bot/pkg/config"
	"github.com/dex-sp/cfg-telegram-bot/pkg/repository"
)

// Singleton implementation.
type Owner struct {
	telegramID   int64
	name         string
	creditCard   string
	personalData config.Owner
	gameList     repository.UserDataRepository
}

func NewOwner(telegramID int64, name string, creditCard string) *Owner {
	return &Owner{telegramID: telegramID, name: name, creditCard: creditCard}
}

func (o *Owner) GetTelegramID() int64 {
	return o.telegramID
}

func (o *Owner) GetName() string {
	return o.name
}

func (o *Owner) GetCreditCard() string {
	return o.creditCard
}
