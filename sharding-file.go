package main

import (
	"flag"
	"io"
	"log"
	"os"
	"time"
)

type fileinfo struct {
	fpath string
	fbyte int
}

func newHandler() *fileinfo {
	return &fileinfo{}
}

func (f *fileinfo) readBlock() {
	FileHandle, err := os.Open(f.fpath)
	if err != nil {
		log.Println(err)
		return
	}
	defer FileHandle.Close()
	buffer := make([]byte, 1024*f.fbyte) // 设置每次读取字节数，
	var count = 0
	for {
		n, err := FileHandle.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		count++
		if n == 0 {
			log.Println("文件读取结束...")
			break
		}
	}
	log.Printf("总共读取次数：%d", count)
}

func main() {
	handler := newHandler()
	flag.StringVar(&handler.fpath, "i", "", "请输入文件")
	flag.IntVar(&handler.fbyte, "b", 1024, "设置每次读取的字节数，默认1024是以M为单位")
	flag.Parse()
	if handler.fpath == "" {
		log.Println("文件路径不能为空，请再次输入！！！")
		os.Exit(0)
	} else if handler.fpath == "./" {
		handler.fpath, _ = os.Getwd()
	}
	now := time.Now()
	defer func() {
		cost := time.Since(now).String()
		log.Printf("总耗时为：%s\n", cost)
	}()
	handler.readBlock()
}
