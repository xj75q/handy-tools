package discover

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

type RegisterInfo struct {
	*EtcdInfo
	Name       string
	IpAddr     string
	NetVersion string
}

func (s *RegisterInfo) Register() {
	var (
		grantLease bool
		leaseId    clientv3.LeaseID
	)
	cli := s.EtcdClient()
	defer cli.Close()

	getRes, err := cli.Get(ctx, s.Name, clientv3.WithCountOnly())
	if err != nil {
		log.Fatalln(err)
	}
	if getRes.Count == 0 {
		grantLease = true
	}
	//租约声明
	if grantLease {
		leaseRes, err := cli.Grant(ctx, 10)
		if err != nil {
			log.Fatalln(err)
		}
		leaseId = leaseRes.ID
	}

	//etcd事务操作
	kv := clientv3.NewKV(cli)
	txn := kv.Txn(ctx)
	_, err = txn.If(clientv3.Compare(clientv3.CreateRevision(s.Name), "=", 0)).
		Then(

			s.leaseInfo(leaseId)...,
		).
		Else(
			s.ignoreLease()...,
		).
		Commit()
	if err != nil {
		log.Fatalln(err)
	}

	//开启永久续约
	if grantLease {
		leaseKeepAlive, err := cli.KeepAlive(ctx, leaseId) //返回的是个Channel
		if err != nil {
			log.Fatalln(err)
		}
		for lease := range leaseKeepAlive {
			fmt.Printf("leaseId :%v, ttl:%d\n", lease.ID, lease.TTL)
		}
	}
}

func (r *RegisterInfo) leaseInfo(leaseId clientv3.LeaseID) []clientv3.Op {
	return []clientv3.Op{
		clientv3.OpPut(r.Name, r.Name, clientv3.WithLease(leaseId)),
		clientv3.OpPut(r.Name+"-ipAddr", r.IpAddr, clientv3.WithLease(leaseId)),
		clientv3.OpPut(r.Name+"-netVersion", r.NetVersion, clientv3.WithLease(leaseId)),
	}
}

func (r *RegisterInfo) ignoreLease() []clientv3.Op {
	return []clientv3.Op{
		clientv3.OpPut(r.Name, r.Name, clientv3.WithIgnoreLease()),
		clientv3.OpPut(r.Name+"-ipAddr", r.IpAddr, clientv3.WithIgnoreLease()),
		clientv3.OpPut(r.Name+"-netVersion", r.NetVersion, clientv3.WithIgnoreLease()),
	}
}
