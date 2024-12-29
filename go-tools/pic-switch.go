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
	"strings"
	"sync"
	"time"
)

var (
	commandName = "convert"
	outPath     string
	picWorkers  = runtime.NumCPU()
)

type fileInfo struct {
	input     string
	output    string
	ptype     string
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

	case ".jpg", ".jpeg", ".jpe", ".jfi":
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
	outType := fmt.Sprintf(".%s", f.ptype)
	fInfo := picStream.(map[string]interface{})
	for pathAndFilename, value := range fInfo {
		info := value.(os.FileInfo)
		picName := fileCommon.GetFileName(info.Name())
		picType := fileCommon.GetFileExt(info.Name())
		input := fileCommon.GetFilePath(pathAndFilename)
		if picType == fmt.Sprintf(".%s", f.ptype) {
			continue
		}
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
		rand.Seed(time.Now().UnixNano())
		if strings.Contains(picName, "!") && picNameLength > 40 {
			var newName string
			content := strings.Split(picName, "!")[0]
			if len(content) > 35 {
				newName = content[32:]
			} else {
				tmpName := content[28:]
				randomNum := rand.Intn(90) + 10 // 生成一个2位随机数
				newName = fmt.Sprintf("%s%d", tmpName, randomNum)
			}
			outFile := fmt.Sprintf("%s%s%s", outPath, newName, outType)
			//log.Println(outFile)
			if err := f.convertCmd(pathAndFilename, outFile); err != nil {
				return err
			}
		} else if picNameLength > 15 && picNameLength < 25 {
			randomNum := rand.Intn(900) + 100 // 生成一个3位随机数
			newName := string([]rune(picName)[:12])
			outName := fmt.Sprintf("%s%d%s", newName, randomNum, outType)
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

		log.Printf("转换成 [%s] 格式完成，源文件 [%s] 将被删除……\n", f.ptype, info.Name())
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
	flag.StringVar(&ph.input, "i", "", "请输入要转换的文件路径")
	flag.StringVar(&ph.output, "o", "./", "转换后的保存路径，默认在当前路径下")
	flag.StringVar(&ph.ptype, "t", "jpg", "请输入要转换为哪种图片格式")
	flag.Parse()
	//防止在任意目录下操作，需输入路径
	if ph.input == "" {
		log.Println("文件路径不能为空，如需在当前目录下，可输入'-i ./'，请继续操作")
		os.Exit(1)
	} else if ph.input == "./" {
		ph.input, _ = os.Getwd()
	}
	outType := fmt.Sprintf(".%s", ph.ptype)
	if !ph.judgeFileType(outType) {
		log.Fatal("请输入正确图片格式...")
		os.Exit(1)
	}

	now := time.Now()
	defer func() {
		cost := time.Since(now)
		costMicro := cost.Microseconds()
		costStr := cost.String()
		if costMicro < 30000 {
			log.Println("此文件夹下没有图片或者都已是目标格式，可查看后继续操作!")
		} else {
			log.Printf("转换完成，总耗时为：%s\n", costStr)
		}
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
