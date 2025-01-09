package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func ConnectRedis(host, port string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})

	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		log.Printf("Redis connection error: %v", err)
		return nil, err
	}

	return rdb, nil
}
