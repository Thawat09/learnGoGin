package database

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
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

	log.Println("Redis connected successfully")

	return rdb, nil
}
