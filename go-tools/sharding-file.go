package main

import (
	"flag"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

type fileinfo struct {
	fpath string
	fbyte int
}

func newHandler() *fileinfo {
	return &fileinfo{}
}

func (f *fileinfo) isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func (f *fileinfo) isFile(path string) bool {
	return !f.isDir(path)
}

func (f *fileinfo) readBlock() {
	if !f.isFile(f.fpath) {
		log.Println(">> 请输入正确的文件及路径...")
		os.Exit(1)
	}

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
	runtime.GOMAXPROCS(runtime.NumCPU())
	handler := newHandler()
	flag.StringVar(&handler.fpath, "i", "", "请输入文件")
	flag.IntVar(&handler.fbyte, "b", 1024, "设置每次读取的字节数，默认1024是以M为单位")
	flag.Parse()
	now := time.Now()
	defer func() {
		cost := time.Since(now).String()
		log.Printf("总耗时为：%s\n", cost)
	}()
	handler.readBlock()
}
