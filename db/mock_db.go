package db

import (
	"context"
	"time"
)

type mockDBHandler struct {
	ctx context.Context
}

func NewMockDBHandler() *mockDBHandler {
	return &mockDBHandler{
		ctx: context.Background(),
	}
}

func (r *mockDBHandler) ConnectToDB(addr string) error {
	return nil
}

func (r *mockDBHandler) SetDailyLimit(limit int) {

}

func (r *mockDBHandler) SetDBSeeds(seeds map[string]string) error {
	return nil
}

func (r *mockDBHandler) GetDBSeeds() (map[string]string, error) {
	return map[string]string{}, nil
}

func (r *mockDBHandler) InsertOrUpdateUUID(uuid string, duration time.Duration) (int64, error) {
	return 0, nil
}
