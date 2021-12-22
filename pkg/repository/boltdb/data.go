package boltdb

import (
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/dex-sp/cfg-telegram-bot/pkg/repository"
)

type UserDataRepository struct {
	db *bolt.DB
}

func NewUserDataRepository(db *bolt.DB) *UserDataRepository {
	return &UserDataRepository{db: db}
}

func (r *UserDataRepository) Save(chatID int64, data string, bucket repository.Bucket) error {

	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(intToBytes(chatID), []byte(data))
	})
}

func (r *UserDataRepository) Get(chatID int64, bucket repository.Bucket) (string, error) {

	var data string

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data = string(b.Get(intToBytes(chatID)))
		return nil
	})

	if err != nil {
		return "", err
	}

	// if data == "" {
	// 	return "", fmt.Errorf("%s data not found", bucket)
	// }
	return data, nil
}

func intToBytes(num int64) []byte {
	return []byte(strconv.FormatInt(num, 10))
}
