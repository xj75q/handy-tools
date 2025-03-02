package discover

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var (
	etcdEndpoints = []string{"127.0.0.1:2379", "127.0.0.1:12379", "127.0.0.1:22379"}
	ctx           = context.Background()
	DialTimeout   = time.Second * 5
)

type EtcdInfo struct{}

func (s *EtcdInfo) EtcdClient() *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: DialTimeout,
	})
	if err != nil {
		panic(err)
	}
	return cli
}
