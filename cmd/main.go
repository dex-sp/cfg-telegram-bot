package main

import (
	"log"
	"os"

	"github.com/dex-sp/cfg-telegram-bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	tokenName string = "TELEGRAM_APITOKEN"
)

func main() {

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Panic(err)
	}

	// Get the TELEGRAM_APITOKEN environment variable
	telegramToken, exists := os.LookupEnv(tokenName)
	if !exists {
		log.Panicf("%s not exists", tokenName)
	}

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	tgBot := telegram.NewBot(bot)
	if err := tgBot.Start(); err != nil {
		log.Fatal(err)
	}
}
