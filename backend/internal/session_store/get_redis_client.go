package session_store

import (
	"github.com/redis/go-redis/v9"

	"context"
)

func GetRedisClient(
	ctx context.Context,
	connConfig UniversalRedisConfig,
) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         connConfig.Hostname + ":" + connConfig.Port,
		Username:     connConfig.Username,
		Password:     connConfig.Password,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
