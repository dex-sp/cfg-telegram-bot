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
	return data, nil
}

func (r *UserDataRepository) Len(bucket repository.Bucket) int64 {

	var sz int64
	r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		sz = int64(b.Stats().KeyN)
		return nil
	})

	return sz
}

func (r *UserDataRepository) GetAll(bucket repository.Bucket) map[int64]string {

	list := make(map[int64]string)

	r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			key, err := strconv.ParseInt(string(k), 10, 64)
			if err != nil {
				return err
			}
			list[key] = string(v)
		}
		return nil
	})
	return list
}

func (r *UserDataRepository) Clear(bucket repository.Bucket) error {

	var err error

	err = r.db.Update(func(tx *bolt.Tx) error {

		err = tx.DeleteBucket([]byte(bucket))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func intToBytes(num int64) []byte {
	return []byte(strconv.FormatInt(num, 10))
}
