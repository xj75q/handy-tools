package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

type netInfo struct {
	Ip      string
	Port    int
	AddChan chan string
	Wg      *sync.WaitGroup
	Workers int
}

func netHandler() *netInfo {
	return &netInfo{
		Wg:      new(sync.WaitGroup),
		AddChan: make(chan string, 500),
		Workers: runtime.NumCPU(),
	}
}

func (n *netInfo) genAddress(ip net.IP) {
	defer n.Wg.Done()
	for port := 1; port <= 65535; port++ {
		address := fmt.Sprintf("%s:%d", ip, port)
		n.AddChan <- address
	}
}

func (n *netInfo) getPort() {
	defer n.Wg.Done()
	for {
		select {
		case address := <-n.AddChan:
			conn, err := net.Dial("tcp", address)
			if err != nil {
				continue
			}
			conn.Close()
			fmt.Println("开放端口为:", address)
		default:
			return
		}
	}
}

func (n *netInfo) parseIP(s string) (net.IP, int) {
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

func main() {
	input := os.Args[1]

	netCtl := netHandler()
	ip, _ := netCtl.parseIP(input)
	if ip == nil {
		fmt.Println(">> 请输入正确的ip地址")
		os.Exit(0)
	}

	fmt.Println(">> 开始扫描")
	defer close(netCtl.AddChan)
	netCtl.Wg.Add(1)
	go netCtl.genAddress(ip)
	time.Sleep(5 * time.Millisecond)
	var now = time.Now()
	defer func() {
		cost := time.Since(now)
		fmt.Printf(">> 总耗时为：%v\n", cost)
	}()
	for i := 1; i < netCtl.Workers; i++ {
		netCtl.Wg.Add(1)
		go netCtl.getPort()
	}

	netCtl.Wg.Wait()
}
