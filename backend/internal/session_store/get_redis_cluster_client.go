package session_store

import (
	"github.com/redis/go-redis/v9"

	"context"
	"errors"
	"net"
	"time"
)

const maxHostLookupRetries = 10
const timeout = time.Second * 10

func GetRedisClusterClient(
	ctx context.Context,
	connConfig UniversalRedisConfig,
) (*redis.ClusterClient, error) {
	_addrs := []string{}
	var err error
	for range maxHostLookupRetries {
		_addrs, err = net.LookupHost(connConfig.Hostname)
		if err == nil {
			break
		}
		time.Sleep(timeout)
	}

	if err != nil {
		return nil, errors.New("Error while looking up ips: " + err.Error())
	}

	addrs := make([]string, len(_addrs))
	for i, addr := range _addrs {
		addrs[i] = addr + ":" + connConfig.Port
	}

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrs,
		Username:     connConfig.Username,
		Password:     connConfig.Password,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	})

	if err = client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
