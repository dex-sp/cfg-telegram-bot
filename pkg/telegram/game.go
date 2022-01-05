package telegram

import (
	"time"

	"github.com/dex-sp/cfg-telegram-bot/pkg/repository"
)

type Game struct {
	date time.Time
	list repository.UserDataRepository
}

func NewGame(list repository.UserDataRepository) *Game {
	return &Game{list: list}
}

func (g *Game) SetDate() error {

	return nil
}

func (g *Game) AddPlayer() error {
	return nil
}

func (g *Game) Gameover() error {
	return nil
}
