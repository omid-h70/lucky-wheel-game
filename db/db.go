package db

import (
	"errors"
	"time"
)

type DBHandler interface {
	ConnectToDB(addr string) error
	InsertOrUpdateUUID(uuid string, duration time.Duration) (int64, error)
	SetDBSeeds(seeds map[string]string) error
	GetDBSeeds() (map[string]string, error)
	SetDailyLimit(limit int)
}

var (
	ErrDBIsUnreachable   = errors.New("DB is Unreachable")
	ErrKeyDoesntExist    = errors.New("Key Doesn't Exist")
	ErrDailyLimitReached = errors.New("Daily Limit Reached")
)
