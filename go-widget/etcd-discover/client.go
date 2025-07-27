package main

import (
	"context"
	"fmt"
	"go-widget/etcd-discover/discover"
	pb "go-widget/etcd-discover/proto-files"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

/*
原先的grpc部分
func getServerAddr() string {
	return "0.0.0.0:50051"
}
*/

func getServerAddr(svcName string) string {
	s := discover.HandlerDiscover.Discover(svcName)
	if s == nil {
		return ""
	}

	return s.IpAddr
}

func sayHello() {
	addr := getServerAddr("discSvc")
	if addr == "" {
		log.Println("未发现可用服务...")
		return
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	log.Println("client start ....")

	client := pb.NewGreeterClient(conn)
	req := &pb.HelloRequest{
		Msg: "client send ...",
	}
	replay, err := client.SayHi(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf(">> %s\n", replay.GetMsg())
}

func main() {
	go discover.HandlerDiscover.WatchSvc("discSvc")
	heartbeat := time.NewTicker(2 * time.Second) // 心跳间隔为2秒
	defer heartbeat.Stop()
	for range heartbeat.C {
		sayHello()
	}
}
