package main

import (
	"context"
	"flag"
	"fmt"
	"go-widget/etcd-discover/discover"
	pb "go-widget/etcd-discover/proto-files"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
)

var (
	ip   = flag.String("ip", "0.0.0.0", "")
	port = flag.Int("port", 50051, "")
)

type (
	SvcRegister struct{}

	GreetServer struct {
		pb.GreeterServer
	}
)

func (h *GreetServer) SayHi(ctx context.Context, in *pb.HelloRequest) (*pb.HelloRespose, error) {
	fmt.Printf(">> recv client msg:%s\n", in.Msg)
	return &pb.HelloRespose{
		Msg: fmt.Sprintf("server resp:%v", in.Msg),
	}, nil
}

func (r *SvcRegister) Register(s grpc.ServiceRegistrar, srv pb.GreeterServer) {
	pb.RegisterGreeterServer(s, srv)
	etcdRegister := &discover.RegisterInfo{
		Name:       "discSvc",
		IpAddr:     fmt.Sprintf("%s:%s", *ip, strconv.Itoa(*port)),
		NetVersion: "grpc",
	}
	go etcdRegister.Register()
}

func main() {
	flag.Parse()
	register := &SvcRegister{}
	s := grpc.NewServer()
	//pb.RegisterGreeterServer(s, &GreetServer{}) //这里是相对于proto来说的改动
	register.Register(s, &GreetServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("server start...")
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
