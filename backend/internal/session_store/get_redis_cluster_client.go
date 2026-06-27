package session_store

import (
	"github.com/redis/go-redis/v9"

	"context"
	"errors"
	"net"
	"time"
)

const max_host_lookup_retires = 10
const timeout = time.Second * 10

func Get_Redis_Cluster_Client(
	ctx context.Context,
	conn_config Universal_Redis_Config,
) (*redis.ClusterClient, error) {
	_addrs := []string{}
	var err error
	for range max_host_lookup_retires {
		_addrs, err = net.LookupHost(conn_config.Hostname)
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
		addrs[i] = addr + ":" + conn_config.Port
	}

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrs,
		Username:     conn_config.Username,
		Password:     conn_config.Password,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	})

	if err = client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
