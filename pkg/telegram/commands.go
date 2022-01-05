package telegram

import (
	"fmt"

	"github.com/dex-sp/cfg-telegram-bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	startCommand    string = "start"
	game0verCommand string = "gameover"
	usersCommand    string = "users"
	playersCommand  string = "players"
)

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {

	msg := tgbotapi.NewMessage(message.Chat.ID, b.config.CommandResponses.Start)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = startMenu
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {

	msg := tgbotapi.NewMessage(message.Chat.ID, b.config.CommandResponses.Default)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleUsersCommand(message *tgbotapi.Message) error {

	isOwner, err := b.isOwner(message.From.ID)
	if err != nil {
		return err
	}

	// Only for owner
	if isOwner {
		text, err := b.getUsersListString(b.userDataRepository, repository.Phones)
		if err != nil {
			return err
		}
		msg := tgbotapi.NewMessage(b.config.Owner.TelegramID, fmt.Sprintf(
			"TOTAL NUMBER OF USERS\n\n%s", text))
		_, err = b.bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return err
}

func (b *Bot) handlePlayersCommand(message *tgbotapi.Message) error {

	isOwner, err := b.isOwner(message.From.ID)
	if err != nil {
		return err
	}

	// Only for owner
	if isOwner {
		text, err := b.getUsersListString(b.userDataRepository, repository.Confirmations)
		if err != nil {
			return err
		}
		msg := tgbotapi.NewMessage(b.config.Owner.TelegramID, fmt.Sprintf(
			"PLAYERS FOR THE NEXT GAME\n\n%s", text))
		_, err = b.bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return err
}

func (b *Bot) handleGame0verCommand(message *tgbotapi.Message) error {

	isOwner, err := b.isOwner(message.From.ID)
	if err != nil {
		return err
	}

	// Only for owner
	if isOwner {
		if b.userDataRepository.Len(repository.Confirmations) > 0 {
			err = b.userDataRepository.Clear(repository.Confirmations)
			if err != nil {
				return err
			}
		}

		msg := tgbotapi.NewMessage(b.config.Owner.TelegramID,
			b.config.CommandResponses.Gameover)
		_, err = b.bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return err
}

func (b *Bot) isOwner(userID int64) (bool, error) {
	var err error
	isOwner := userID != b.config.Owner.TelegramID

	if isOwner {
		msg := tgbotapi.NewMessage(userID,
			b.config.Errors.NotEnoughRights)
		_, err = b.bot.Send(msg)
		if err != nil {
			return isOwner, err
		}
	}
	return isOwner, err
}
