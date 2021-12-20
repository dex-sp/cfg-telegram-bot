package main

import (
	"log"
	"os"

	"github.com/dex-sp/cfg-telegram-bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	telegramTokenVarName string = "TELEGRAM_APITOKEN"
	// trelloAppKeyVarName  string = "TRELLO_APP_KEY"
)

func main() {

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// Get the TELEGRAM_APITOKEN environment variable
	telegramToken, exists := os.LookupEnv(telegramTokenVarName)
	if !exists {
		log.Fatalf("%s not exists", telegramTokenVarName)
	}

	// Get the TRELLO_APP_KEY environment variable
	// trelloAppKey, exists := os.LookupEnv(trelloAppKeyVarName)
	// if !exists {
	// 	log.Panicf("%s not exists", trelloAppKeyVarName)
	// }

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	tgBot := telegram.NewBot(bot)
	if err := tgBot.Start(); err != nil {
		log.Fatal(err)
	}
}
