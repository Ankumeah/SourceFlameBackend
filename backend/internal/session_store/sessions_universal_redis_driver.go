package session_store

import (
	"github.com/redis/go-redis/v9"

	"context"
	"errors"
	"strconv"
	"time"
)

type RedisClient interface {
	ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd
	ZScore(ctx context.Context, key string, member string) *redis.FloatCmd
	ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	ZCount(ctx context.Context, key string, min string, max string) *redis.IntCmd
}
type sessionsUniversalRedisDriver struct{ rdb RedisClient }

func GetSessionsUniversalRedisDriver(
	client RedisClient,
	tokenTimeout time.Duration,
	tokenLimit int,
	tokenLength int,
	tokenNamespace string,
) *SessionStore {
	return &SessionStore{
		&sessionsUniversalRedisDriver{client},
		tokenTimeout,
		tokenLimit,
		tokenLength,
		tokenNamespace,
	}
}

func (r *sessionsUniversalRedisDriver) AddSession(ctx context.Context, username string, token string, timeout time.Duration) error {
	now := time.Now().Add(timeout).Unix()
	return r.rdb.ZAdd(ctx, username, redis.Z{Member: token, Score: float64(now)}).Err()
}

func (r *sessionsUniversalRedisDriver) ValidateSession(ctx context.Context, username string, token string) (bool, error) {
	now := float64(time.Now().Unix())
	exp, err := r.rdb.ZScore(ctx, username, token).Result()

	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	} else if exp <= now {
		return false, nil
	} else {
		return true, nil
	}
}

func (r *sessionsUniversalRedisDriver) DeleteSession(ctx context.Context, username string, token string) error {
	return r.rdb.ZRem(ctx, username, token).Err()
}

func (r *sessionsUniversalRedisDriver) GetSessionCount(ctx context.Context, username string) (int, error) {
	now := time.Now().Unix()
	count, err := r.rdb.ZCount(ctx, username, strconv.FormatInt(now, 10), "inf").Result()
	if err != nil {
		return 0, err
	} else {
		return int(count), nil
	}
}
