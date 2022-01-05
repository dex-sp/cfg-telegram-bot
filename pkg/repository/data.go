package repository

type Bucket string

const (
	Phones        Bucket = "phones"
	Locations     Bucket = "locations"
	Confirmations Bucket = "confirmations"
)

type UserDataRepository interface {
	Save(chatID int64, data string, bucket Bucket) error
	Get(chatID int64, bucket Bucket) (string, error)
	Len(bucket Bucket) int64
	GetAll(bucket Bucket) map[int64]string
	Clear(bucket Bucket) error
}
