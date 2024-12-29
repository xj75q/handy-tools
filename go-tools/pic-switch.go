package main

import (
	"flag"
	"fmt"
	"go-tools/fileCommon"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	commandName = "convert"
	outType     = ".jpg"
	outPath     string
	picWorkers  = runtime.NumCPU()
)

type fileInfo struct {
	input     string
	output    string
	wg        *sync.WaitGroup
	eventChan chan interface{}
	exitChan  chan bool
}

func handlerPic() *fileInfo {
	return &fileInfo{
		wg:        new(sync.WaitGroup),
		eventChan: make(chan interface{}, 20),
		exitChan:  make(chan bool, 1),
	}
}

func (f *fileInfo) judgeFileType(flag string) bool {
	switch flag {
	case ".webp", ".awebp":
		return true

	case ".png", ".bmg", ".gif":
		return true

	case ".jpeg", ".jpe", ".jfi":
		return true

	case ".avif", ".mng", ".jng":
		return true

	case ".tga", ".wmf", ".dng", ".pnm", ".pgm", ".ppm":
		return true

	default:
		return false
	}
}

func (f *fileInfo) handleFolder() error {
	err := filepath.Walk(f.input, func(pathAndFilename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fileExt := fileCommon.GetFileExt(info.Name())
		isPicType := f.judgeFileType(fileExt)
		if fileCommon.IsFile(pathAndFilename) && isPicType {
			var fInfo = make(map[string]interface{})
			fInfo[pathAndFilename] = info
			f.wg.Add(1)
			go func() {
				defer f.wg.Done()
				f.eventChan <- fInfo
			}()
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (f *fileInfo) handleEventStream() {
	defer f.wg.Done()
	for {
		select {
		case picStream := <-f.eventChan:
			if err := f.convertPic(picStream); err != nil {
				log.Printf(">> 转换图片出错：%v", err)
				return
			}
		default:
			return
		}
	}
}

func (f *fileInfo) convertPic(picStream interface{}) error {
	fInfo := picStream.(map[string]interface{})
	for pathAndFilename, value := range fInfo {
		info := value.(os.FileInfo)
		picName := fileCommon.GetFileName(info.Name())
		input := fileCommon.GetFilePath(pathAndFilename)
		picNameLength := len([]rune(picName))
		if f.output == "./" {
			outPath = fmt.Sprintf("%s%s", input, string(os.PathSeparator))
		} else {
			err, tmpOut := f.outPutOperate()
			if err != nil {
				return err
			}
			outPath = fmt.Sprintf("%s%s", tmpOut, string(os.PathSeparator))
		}
		if strings.Contains(picName, "!") && picNameLength > 40 {
			content := strings.Split(picName, "!")[0]
			newName := content[32:]
			outName := fmt.Sprintf("%s%s", newName, outType)
			outFile := fmt.Sprintf("%s%s", outPath, outName)
			//log.Println(outFile)
			if err := f.convertCmd(pathAndFilename, outFile); err != nil {
				return err
			}
		} else if picNameLength > 15 && picNameLength < 25 {
			rand.Seed(time.Now().UnixNano())
			randomNum := rand.Intn(900) + 100 // 生成一个3位随机数
			newName := string([]rune(picName)[:12])
			outName := fmt.Sprintf("%s%s%s", newName, strconv.Itoa(randomNum), outType)
			outFile := fmt.Sprintf("%s%s", outPath, outName)
			//log.Println(newName, outFile)
			if err := f.convertCmd(pathAndFilename, outFile); err != nil {
				return err
			}
		} else {
			outFile := fmt.Sprintf("%s%s%s", outPath, picName, outType)
			//log.Println(outFile)
			if err := f.convertCmd(pathAndFilename, outFile); err != nil {
				return err
			}
		}
		log.Printf("转换完成，源文件 [%s] 将被删除……\n", info.Name())
		//time.Sleep(100 * time.Millisecond) //为了方便看到过程，实际可以删除
		os.Remove(pathAndFilename)
	}
	return nil
}

func (f *fileInfo) convertCmd(pathAndFilename, outFile string) error {
	cmd := exec.Command(commandName, pathAndFilename, outFile)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func (f *fileInfo) outPutOperate() (error, string) {
	out := filepath.Clean(f.output)
	if fileCommon.CheckSavePath(out) {
		if err := fileCommon.CreateSavePath(out, os.FileMode(0755)); err != nil {
			return fmt.Errorf("创建目录失败![%v]\n", err), ""
		}
	}
	permission := fileCommon.CheckPermission(out)
	if !permission {
		if err := os.Chmod(out, os.FileMode(0755)); err != nil {
			return fmt.Errorf("目录赋权限失败![%v]\n", err), ""
		}
	}
	return nil, out
}

func main() {
	ph := handlerPic()
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&ph.input, "i", "", "请输入路径")
	flag.StringVar(&ph.output, "o", "./", "请输入路径")
	flag.Parse()
	if ph.input == "" {
		log.Println("文件路径不能为空，请再次输入！！！")
		os.Exit(1)
	} else if ph.input == "./" {
		ph.input, _ = os.Getwd()
	}
	now := time.Now()
	defer func() {
		cost := time.Since(now).String()
		log.Printf("总耗时为：%s\n", cost)
	}()
	defer close(ph.eventChan)
	if err := ph.handleFolder(); err != nil {
		log.Printf(">> 处理文件夹时出错: %v", err)
		return
	}
	ph.wg.Add(1)
	go func() {
		defer ph.wg.Done()
		for i := 1; i < picWorkers; i++ {
			ph.wg.Add(1)
			go ph.handleEventStream()
		}
	}()
	ph.wg.Wait()
}
