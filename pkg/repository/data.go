package repository

type Bucket string

const (
	Phones    Bucket = "phones"
	Locations Bucket = "locations"
)

type UserDataRepository interface {
	Save(chatID int64, data string, bucket Bucket) error
	Get(chatID int64, bucket Bucket) (string, error)
}
