package main

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/dex-sp/cfg-telegram-bot/pkg/config"
	"github.com/dex-sp/cfg-telegram-bot/pkg/repository"
	"github.com/dex-sp/cfg-telegram-bot/pkg/repository/boltdb"
	"github.com/dex-sp/cfg-telegram-bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true

	userDB, err := initDB(cfg.UserDBPath)
	if err != nil {
		log.Fatal(err)
	}

	userDataRepository := boltdb.NewUserDataRepository(userDB)

	tgBot := telegram.NewBot(bot, userDataRepository, cfg)
	if err := tgBot.Start(); err != nil {
		log.Fatal(err)
	}
}

func initDB(dbPath string) (*bolt.DB, error) {

	db, err := bolt.Open(dbPath, 0600, nil)
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
		_, err = tx.CreateBucketIfNotExists([]byte(repository.Confirmations))
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
