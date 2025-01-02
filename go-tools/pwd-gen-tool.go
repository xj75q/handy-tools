package main

import (
	"flag"
	"fmt"
	"go-tools/fileCommon"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"time"
)

var (
	strSet                = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()_+`-=[]\\{}|;':\",./<>? ")
	markpwd        uint32 = 0
	fileStore      string
	currentPath, _ = os.Getwd()
)

type Param struct {
	SavePath  string `json:"savepath"`
	Pwdlength int    `json:"pwdlen"`
	PwdStart  int
}

func paramHandler() *Param {
	return &Param{}
}

func (p *Param) genNewStr(strSet string, n int, sc chan string) {
	for i, c := range strSet {
		if n == 1 {
			sc <- string(c)
		} else {
			var ssc = make(chan string)
			go p.genNewStr(strSet[:i]+strSet[i+1:], n-1, ssc)
			for k := range ssc {
				sc <- fmt.Sprintf("%v%v", string(c), k)
			}
		}
	}
	close(sc)
}

func (p *Param) execGen(start, length int, fs *os.File) {
	if start > length {
		log.Println("密码起始位不得大于密码长度")
		return
	}
	for i := start; i <= length; i++ {
		log.Println("i:", i)
		sc := make(chan string)
		go p.genNewStr(string(strSet), i, sc)
		for x := range sc {
			atomic.AddUint32(&markpwd, 1)
			fs.WriteString(x)
			fs.WriteString(string("\n"))
			log.Println("生成新的pwd为:", x)
		}
	}
}

func (p *Param) execPwd() error {
	starttime := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())
	newName := fmt.Sprintf("%s%d%s", "password_", p.Pwdlength, "位.txt")
	if p.SavePath == "" {
		fileStore = fmt.Sprintf("%s%s%s", currentPath, fileCommon.PathMark, newName)
	} else {
		fileStore = fmt.Sprintf("%s%s%s", filepath.Clean(p.SavePath), fileCommon.PathMark, newName)
	}

	if fileCommon.CheckSavePath(fileStore) {
		if err := fileCommon.CreateFile(fileStore); err != nil {
			return err
		}
	}

	fs, err := os.OpenFile(fileStore, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer fs.Close()
	p.execGen(p.PwdStart, p.Pwdlength, fs)
	imarkFinal := atomic.LoadUint32(&markpwd)
	since := int(time.Since(starttime).Seconds())
	log.Println("完成消耗时间:", since, "s", "生成:", imarkFinal, "个密码")
	time.Sleep(10)
	return nil
}

func main() {
	inputParam := paramHandler()
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&inputParam.SavePath, "o", "", "文件保存路径")
	flag.IntVar(&inputParam.Pwdlength, "l", 3, "生成的密码长度")
	flag.IntVar(&inputParam.PwdStart, "s", 1, "只要几位数开始的密码，默认从1位开始")
	flag.Parse()

	err := inputParam.execPwd()
	if err != nil {
		log.Printf(">> err: %v", err)
	}
}
