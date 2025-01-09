package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ctx         = context.Background()
	RedisClient *redis.Client
)

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

func SetRedisClient(client *redis.Client) {
	RedisClient = client
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

func SetValue(rdb *redis.Client, key, value string, expiration time.Duration) error {
	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Printf("Failed to set key %s in Redis: %v", key, err)
		return err
	}

	log.Printf("Key %s set successfully in Redis", key)
	return nil
}

func GetValue(rdb *redis.Client, key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("Key %s does not exist in Redis", key)
		return "", nil
	} else if err != nil {
		log.Printf("Failed to get key %s from Redis: %v", key, err)
		return "", err
	}

	return val, nil
}

func DeleteKey(rdb *redis.Client, key string) error {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Failed to delete key %s from Redis: %v", key, err)
		return err
	}

	log.Printf("Key %s deleted successfully from Redis", key)
	return nil
}

func KeyExists(rdb *redis.Client, key string) (bool, error) {
	val, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		log.Printf("Failed to check existence of key %s in Redis: %v", key, err)
		return false, err
	}

	return val > 0, nil
}
