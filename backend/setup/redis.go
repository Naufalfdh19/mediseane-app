package setup

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func ConnectRedisDB(c context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(c).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
