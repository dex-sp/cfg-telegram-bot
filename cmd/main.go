package main

import (
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/dex-sp/cfg-telegram-bot/pkg/repository"
	"github.com/dex-sp/cfg-telegram-bot/pkg/repository/boltdb"
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

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	userDataRepository := boltdb.NewUserDataRepository(db)

	tgBot := telegram.NewBot(bot, userDataRepository)
	if err := tgBot.Start(); err != nil {
		log.Fatal(err)
	}
}

func initDB() (*bolt.DB, error) {

	db, err := bolt.Open("userData.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.Phones))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(repository.Locations))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return db, nil
}
