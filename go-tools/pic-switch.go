package main

import (
	"flag"
	"fmt"
	"go-tools/fileCommon"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	commandName = "convert"
	outType     = ".jpg"
	picWorkers  = runtime.NumCPU()
)

type fileInfo struct {
	fPath     string
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
	case "webp", "awebp":
		return true

	case "png", "bmg", "gif":
		return true

	case "jpeg", "jpe", "jfi":
		return true

	case "avif", "mng", "jng":
		return true

	case "tga", "wmf", "dng", "pnm", "pgm", "ppm":
		return true

	default:
		return false
	}
}

func (f *fileInfo) executeSwitch() {

	err := filepath.Walk(f.fPath, func(pathAndFilename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		infoName := strings.Split(strings.ToLower(info.Name()), ".")
		flag := infoName[len(infoName)-1]
		isPicType := f.judgeFileType(flag)
		if fileCommon.IsFile(pathAndFilename) && isPicType {
			var fInfo = make(map[string]interface{})
			fInfo[pathAndFilename] = info
			f.wg.Add(1)
			go func() {
				defer f.wg.Done()
				f.eventChan <- fInfo
			}()
			return nil
		}
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}

}

func (f *fileInfo) handleEventStream() {
	defer f.wg.Done()
	for {
		select {
		case picStream := <-f.eventChan:
			f.convertPic(picStream)
		default:
			return
		}
	}
}

func (f *fileInfo) convertPic(picStream interface{}) {
	fInfo := picStream.(map[string]interface{})
	for pathAndFilename, value := range fInfo {
		info := value.(os.FileInfo)
		infoName := strings.Split(strings.ToLower(info.Name()), ".")
		picName := strings.Join(infoName[:len(infoName)-1], ".")
		picNameLength := len([]rune(picName))
		outlist := strings.Split(pathAndFilename, string(os.PathSeparator))
		final := outlist[:len(outlist)-1]
		outpath := strings.Join(final, "/") + string(os.PathSeparator)
		if strings.Contains(picName, "!") && picNameLength > 40 {
			content := strings.Split(picName, "!")[0]
			outName := content[32:] + outType
			outContent := outpath + outName
			cmd := exec.Command(commandName, pathAndFilename, outContent)
			if err := cmd.Start(); err != nil {
				fmt.Println(err.Error())
			}
			if err := cmd.Wait(); err != nil {
				fmt.Println(err.Error())
			}
		} else if picNameLength > 15 && picNameLength < 25 {
			//todo 名字去重
			outName := string([]rune(picName)[:15]) + outType
			outContent := outpath + outName
			fmt.Println(outContent)
			cmd := exec.Command(commandName, pathAndFilename, outContent)
			if err := cmd.Start(); err != nil {
				fmt.Println(err.Error())
			}
			if err := cmd.Wait(); err != nil {
				fmt.Println(err.Error())
			}
		} else {
			outContent := fmt.Sprintf("%s%s%s", outpath, picName, outType)
			//fmt.Println(outContent)
			cmd := exec.Command(commandName, pathAndFilename, outContent)
			if err := cmd.Start(); err != nil {
				fmt.Println(err.Error())
			}
			if err := cmd.Wait(); err != nil {
				fmt.Println(err.Error())
			}
		}
		fmt.Printf("转换完成，源文件 [%s] 将被删除……\n", info.Name())
		time.Sleep(1 * time.Second)
		os.Remove(pathAndFilename)
	}
}

type filterSameName struct{}

func (f filterSameName) unique(params []string) []string {
	names := make([]string, 0)
	tmpName := make(map[string]byte)
	for _, v := range params {
		if _, ok := tmpName[v]; !ok {
			tmpName[v] = 1
			names = append(names, v)
		}
	}
	return names
}

func (f filterSameName) fileName(names []string) []string {
	length := len(names)
	if length > 1 {
		return names
	}

	return nil
}

//func

func main() {
	ph := handlerPic()
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&ph.fPath, "i", "", "请输入路径")
	flag.Parse()
	if ph.fPath == "" {
		fmt.Println("文件路径不能为空，请再次输入！！！")
		os.Exit(0)
	} else if ph.fPath == "./" {
		ph.fPath, _ = os.Getwd()
	}
	now := time.Now()
	defer func() {
		cost := time.Since(now).String()
		fmt.Printf("总耗时为：%s\n", cost)
	}()
	defer close(ph.eventChan)
	ph.executeSwitch()

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
