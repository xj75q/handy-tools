package discover

import (
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
)

var HandlerDiscover = &discoverService{
	svc: map[string]*RegisterInfo{},
}

type discoverService struct {
	*EtcdInfo
	sync.RWMutex
	svc map[string]*RegisterInfo //serverName 为key值，信息为value的
}

func (d *discoverService) Discover(svcName string) *RegisterInfo {
	var s *RegisterInfo = nil
	HandlerDiscover.RLock()
	defer HandlerDiscover.RUnlock()
	s, _ = HandlerDiscover.svc[svcName]
	return s
}

func (d *discoverService) WatchSvc(svcName string) {
	cli := d.EtcdClient()
	defer cli.Close()
	getRes, err := cli.Get(ctx, svcName, clientv3.WithPrefix())
	if err != nil {
		log.Fatalln(err)
	}
	if getRes.Count > 0 {
		mp := sliceToMap(getRes.Kvs)
		s := &RegisterInfo{}
		if kv, ok := mp[svcName]; ok {
			s.Name = string(kv.Value)
		}

		if kv, ok := mp[svcName+"-ipAddr"]; ok {
			s.IpAddr = string(kv.Value)
		}

		if kv, ok := mp[svcName+"-netVersion"]; ok {
			s.NetVersion = string(kv.Value)
		}
		HandlerDiscover.Lock()
		HandlerDiscover.svc[svcName] = s
		HandlerDiscover.Unlock()
	}

	d.etcdWatch(cli, svcName)
}

func (d *discoverService) etcdWatch(cli *clientv3.Client, svcName string) {
	rch := cli.Watch(ctx, svcName, clientv3.WithPrefix())
	for wres := range rch {
		for _, ev := range wres.Events {
			if ev.Type == clientv3.EventTypeDelete {
				HandlerDiscover.Lock()
				defer HandlerDiscover.Unlock()
				delete(HandlerDiscover.svc, svcName)
			}

			if ev.Type == clientv3.EventTypePut {
				HandlerDiscover.Lock()
				if _, ok := HandlerDiscover.svc[svcName]; !ok {
					HandlerDiscover.svc[svcName] = &RegisterInfo{} //需要写进去
				}
				defer HandlerDiscover.Unlock()
				switch string(ev.Kv.Key) {
				case svcName:
					HandlerDiscover.svc[svcName].Name = string(ev.Kv.Value)
				case svcName + "-ipAddr":
					HandlerDiscover.svc[svcName].IpAddr = string(ev.Kv.Value)
				case svcName + "-netVersion":
					HandlerDiscover.svc[svcName].NetVersion = string(ev.Kv.Value)
				}
			}
		}
	}
}

func sliceToMap(list []*mvccpb.KeyValue) map[string]*mvccpb.KeyValue {
	mp := make(map[string]*mvccpb.KeyValue, 0)
	for _, item := range list {
		mp[string(item.Key)] = item
	}
	return mp
}
