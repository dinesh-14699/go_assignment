package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var ctx = context.Background()

func InitializeCache(redisAddr string, redisPassword string, redisDB int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,      
		Password: redisPassword,  
		DB:       redisDB,        
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	} else {
		fmt.Println("Successfully connected to Redis!")
	}
}

func SetValue(key string, value string, expiry int) error {
	if expiry > 0 {
		err := rdb.Set(ctx, key, value, time.Duration(expiry)*time.Second).Err() 
		if err != nil {
			return fmt.Errorf("could not set value with expiry: %v", err)
		}
		fmt.Printf("Set key %s with expiry time of %d seconds\n", key, expiry)
	} else {
		err := rdb.Set(ctx, key, value, 0).Err()
		if err != nil {
			return fmt.Errorf("could not set value: %v", err)
		}
		fmt.Printf("Set value: %s = %s\n", key, value)
	}
	return nil
}


func GetValue(key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key does not exist")
		}
		return "", fmt.Errorf("could not get value: %v", err)
	}
	return val, nil
}

func DeleteKey(key string) error {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("could not delete key: %v", err)
	}
	fmt.Printf("Deleted key: %s\n", key)
	return nil
}
