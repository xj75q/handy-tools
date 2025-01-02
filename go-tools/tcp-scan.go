package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

type (
	NetScan interface {
		GenAddress(ip net.IP)
		ParseIP(s string) (net.IP, int)
		GetPort()
	}

	netInfo struct {
		Ip      string
		Port    int
		AddChan chan string
		Wg      *sync.WaitGroup
		Workers int
	}
)

func netHandler() *netInfo {
	return &netInfo{
		Wg:      new(sync.WaitGroup),
		AddChan: make(chan string, 500),
		Workers: runtime.NumCPU(),
	}
}

func (n *netInfo) ParseIP(s string) (net.IP, int) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, 0
	}
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return ip, 4
		case ':':
			return ip, 6
		}
	}
	return nil, 0
}

func (n *netInfo) GenAddress(ip net.IP) {
	defer n.Wg.Done()
	for port := 1; port <= 65535; port++ {
		address := fmt.Sprintf("%s:%d", ip, port)
		n.AddChan <- address
	}
}

func (n *netInfo) GetPort() {
	defer n.Wg.Done()
	for {
		select {
		case address := <-n.AddChan:
			conn, err := net.Dial("tcp", address)
			if err != nil {
				continue
			}
			conn.Close()
			log.Println("开放端口为:", address)
		default:
			return
		}
	}
}

// 使用os.args来获取命令行参数，也可以改为flag
// ./tcp-scan  192.168.49.1
func main() {
	in := os.Args
	if len(in) < 2 {
		log.Println("请输入要扫描的地址")
		os.Exit(1)
	}
	input := os.Args[1]
	netCtl := netHandler()
	ip, _ := netCtl.ParseIP(input)
	if ip == nil {
		log.Println(">> 请输入正确的ip地址")
		os.Exit(1)
	}

	log.Println(">> 开始扫描....")
	defer close(netCtl.AddChan)

	netCtl.Wg.Add(1)
	go netCtl.GenAddress(ip)
	time.Sleep(5 * time.Millisecond)
	var now = time.Now()
	defer func() {
		cost := time.Since(now).String()
		log.Printf(">> 总耗时为：%v\n", cost)
	}()

	for i := 1; i < netCtl.Workers; i++ {
		netCtl.Wg.Add(1)
		go netCtl.GetPort()
	}
	netCtl.Wg.Wait()
}
