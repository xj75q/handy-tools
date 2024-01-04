package main

import (
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"
)

var (
	ip = "192.168.155.122"
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

func (n *netInfo) genAddress() {
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

func main() {
	fmt.Println(">> 开始扫描")
	netCtl := netHandler()
	defer close(netCtl.AddChan)
	netCtl.Wg.Add(1)
	go netCtl.genAddress()
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
