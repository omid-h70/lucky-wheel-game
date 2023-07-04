package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

const (
	SEEDS_KEY           = "seeds#"
	USER_KEY            = "user#"
	DEFAULT_RETRY       = "0"
	DEFAULT_RETRY_LIMIT = "3"
)

type redisHandler struct {
	rdb        *redis.Client
	ctx        context.Context
	dailyLimit int
}

func NewRedisHandler() *redisHandler {
	return &redisHandler{
		ctx:        context.Background(),
		dailyLimit: 3,
	}
}

func (r *redisHandler) ConnectToDB(addr string) error {
	var isConnected bool

	for i := 0; i < 100; i++ {
		r.rdb = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		statusCode := r.rdb.Ping(context.Background())
		fmt.Println("status ", statusCode.String())

		if statusCode.String() == "ping: PONG" {
			//"Connection Not Opened"
			fmt.Printf("Connected On %d Attempt \n", i)
			isConnected = true
			break
		} else {
			fmt.Println("Not Ready Yet...")
		}
		time.Sleep(2 * time.Second)
	}

	if !isConnected {
		log.Fatal("Redis DB IS Down")
		return ErrDBIsUnreachable
	}
	return nil
}

func (r *redisHandler) SetDailyLimit(limit int) {
	r.dailyLimit = limit
}

func (r *redisHandler) SetDBSeeds(seeds map[string]string) error {
	return r.rdb.HSet(r.ctx, "seeds#", seeds).Err()
	/*
		if _, err := r.rdb.Pipelined(r.ctx, func(rdb redis.Pipeliner) error {
			r.rdb.HSet(r.ctx, SEEDS_KEY, seeds)
			r.rdb.Expire(r.ctx, SEEDS_KEY, 0)
			return nil
		}); err != nil {
			return err
		}
		return nil
	*/
}

func (r *redisHandler) GetDBSeeds() (map[string]string, error) {
	return r.rdb.HGetAll(r.ctx, SEEDS_KEY).Result()
}

func (r *redisHandler) InsertOrUpdateUUID(uuid string, duration time.Duration) (int64, error) {
	var cnt int64
	val, err := r.rdb.Get(r.ctx, USER_KEY+uuid).Result()
	if err != nil && err != redis.Nil {
		return 0, ErrDBIsUnreachable
	}

	limit, err := strconv.Atoi(val)

	if len(val) == 0 {
		err = r.rdb.Set(r.ctx, USER_KEY+uuid, DEFAULT_RETRY, duration).Err()
	} else if limit < r.dailyLimit {
		result := r.rdb.Incr(r.ctx, USER_KEY+uuid)
		err, cnt = result.Err(), result.Val()
	} else {
		err, cnt = ErrDailyLimitReached, int64(r.dailyLimit)
	}

	return cnt, err
}
