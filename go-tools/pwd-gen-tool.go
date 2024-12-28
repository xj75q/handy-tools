package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	strSet            = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()_+`-=[]\\{}|;':\",./<>? ")
	exepath, _        = os.Executable()
	current    string = filepath.Dir(exepath)
	imarkpwd   uint32 = 0
)

type folder struct{}

func folderHandler() *folder {
	return &folder{}
}

func (fd *folder) PathExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (fd *folder) CreateFile(InputPath string) string {
	passwordPath := InputPath + string(os.PathSeparator) + "password_file"
	isExists := fd.PathExists(passwordPath)
	if !isExists {
		err := os.Mkdir(passwordPath, 0777)
		if err != nil {
			log.Println(err.Error())
		}
		log.Println("===存放密码文件夹创建成功...")
	}
	return passwordPath
}

func genNewStr(strSet string, n int, sc chan string) {
	for i, c := range strSet {
		if n == 1 {
			sc <- string(c)
		} else {
			var ssc = make(chan string)
			go genNewStr(strSet[:i]+strSet[i+1:], n-1, ssc)
			for k := range ssc {
				sc <- fmt.Sprintf("%v%v", string(c), k)
			}
		}
	}
	close(sc)
}

func execGen(start, length int, fs *os.File) {
	if start > length {
		log.Println("密码起始位不得大于密码长度")
		return
	}
	for i := start; i <= length; i++ {
		log.Println("i:", i)
		sc := make(chan string)
		go genNewStr(string(strSet), i, sc)
		for x := range sc {
			atomic.AddUint32(&imarkpwd, 1)
			fs.WriteString(x)
			fs.WriteString(string("\n"))
			log.Println("生成新的pwd为:", x)
		}
	}
}

func (p *Param) run() error {
	starttime := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())
	fileName := "password_" + strconv.Itoa(p.Pwdlength) + "位.txt"
	fi := current + string(os.PathSeparator) + fileName
	file := filepath.Clean(fi)
	if p.SavePath == "./password.txt" {
		fs, e := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if e != nil {
			return e
		}
		defer fs.Close()
		execGen(p.PwdStart, p.Pwdlength, fs)
	} else {
		if _, err := os.Stat(p.SavePath); err != nil {
			return err
		}
		fd := folderHandler()
		pwd_dir := fd.CreateFile(p.SavePath)
		filename := pwd_dir + string(os.PathSeparator) + fileName
		log.Println(filename)
		fs, e := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if e != nil {
			return e
		}
		defer fs.Close()
		execGen(p.PwdStart, p.Pwdlength, fs)
	}

	imarkFinal := atomic.LoadUint32(&imarkpwd)
	since := int(time.Since(starttime).Seconds())
	log.Println("完成消耗时间:", since, "s", "生成:", imarkFinal, "个密码")
	time.Sleep(10)
	return nil
}

type Param struct {
	SavePath  string `json:"savepath"`
	Pwdlength int    `json:"pwdlen"`
	PwdStart  int
}

func paramHandler() *Param {
	return &Param{}
}

func main() {
	inputParam := paramHandler()
	flag.StringVar(&inputParam.SavePath, "o", "./password.txt", "文件保存路径")
	flag.IntVar(&inputParam.Pwdlength, "l", 3, "生成的密码长度")
	flag.IntVar(&inputParam.PwdStart, "s", 1, "只要几位数开始的密码，默认从1位开始")
	flag.Parse()
	err := inputParam.run()
	if err != nil {
		log.Printf(">> err: %v", err)
	}
}
