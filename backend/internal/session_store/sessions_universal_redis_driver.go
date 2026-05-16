package session_store

import (
	"github.com/redis/go-redis/v9"

	"context"
	"time"
  "strconv"
)

type Redis_client interface {
  ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd
  ZScore(ctx context.Context, key string, member string) *redis.FloatCmd
  ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
  ZCount(ctx context.Context, key string, min string, max string) *redis.IntCmd
}
type sessions_uinversal_redis_driver struct { rdb Redis_client }
func Get_Sessions_Uinversal_Redis_Driver(client Redis_client) *Session_store {
  return &Session_store {
    &sessions_uinversal_redis_driver { client },
  }
}

func (r *sessions_uinversal_redis_driver) Add_Session(ctx context.Context, username string, token string, timeout time.Duration) error {
  now := time.Now().Add(timeout).Unix()
  return r.rdb.ZAdd(ctx, username, redis.Z { Member: token, Score: float64(now) }).Err()
}

func (r *sessions_uinversal_redis_driver) Validate_Session(ctx context.Context, username string, token string) (bool, error) {
  now := float64(time.Now().Unix())
  exp, err := r.rdb.ZScore(ctx, username, token).Result()

  if err == redis.Nil {
    return false, nil
  } else if err != nil {
    return false, err
  } else if exp <= now {
    return false, nil
  } else {
    return true, nil
  }
}

func (r *sessions_uinversal_redis_driver) Delete_Session(ctx context.Context, username string, token string) error {
  return r.rdb.ZRem(ctx, username, token).Err()
}

func (r *sessions_uinversal_redis_driver) Get_Session_Count(ctx context.Context, username string) (uint8, error) {
  now := time.Now().Unix()
  count , err := r.rdb.ZCount(ctx, username, strconv.FormatInt(now, 10), "inf").Result()
  if err != nil {
    return 0, err
  } else {
    return uint8(count), nil
  }
}
