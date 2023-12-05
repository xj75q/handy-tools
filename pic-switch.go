package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	commandName = "convert"
	outType     = ".jpg"
)

type fileInfo struct {
	fPath string
	fType string
}

func handlerPic() *fileInfo {
	return &fileInfo{}
}

func (f *fileInfo) judgeFileType() string {
	switch f.fType {
	case "webp":
		return ".webp"
	case "png":
		return ".png"
	case "jpeg":
		return ".jpeg"
	default:
		return ".webp"
	}
}

func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func isFile(path string) bool {
	return !isDir(path)
}

func (f *fileInfo) executeSwitch() error {
	suffixFlag := f.judgeFileType()
	err := filepath.Walk(f.fPath, func(pathAndFilename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if isFile(pathAndFilename) && strings.HasSuffix(strings.ToLower(info.Name()), suffixFlag) {
			picName := strings.Split(info.Name(), suffixFlag)[0]
			outlist := strings.Split(pathAndFilename, string(os.PathSeparator))
			final := outlist[:len(outlist)-1]
			outpath := strings.Join(final, "/") + string(os.PathSeparator)
			if strings.Contains(picName, "!") && len(picName) > 40 {
				content := strings.Split(picName, "!")[0]
				outName := content[32:] + outType
				outContent := outpath + outName
				cmd := exec.Command(commandName, pathAndFilename, outContent)
				err := cmd.Run()
				if err != nil {
					return err
				}
			} else if len(picName) > 15 && len(picName) < 25 {
				outName := picName[15:] + outType
				outContent := outpath + string(os.PathSeparator) + outName
				cmd := exec.Command(commandName, pathAndFilename, outContent)
				err := cmd.Run()
				if err != nil {
					return err
				}
			} else {
				outContent := outpath + string(os.PathSeparator) + info.Name()
				cmd := exec.Command(commandName, pathAndFilename, outContent)
				err := cmd.Run()
				if err != nil {
					return err
				}
			}
			fmt.Printf("转换完成，源文件 [%s] 将被删除……\n", info.Name())
			if err := os.Remove(pathAndFilename); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	ph := handlerPic()
	flag.StringVar(&ph.fPath, "i", "", "请输入路径")
	flag.Parse()
	if ph.fPath == "" {
		panic("文件路径不能为空，请再次输入！！！")
	} else if ph.fPath == "./" {
		ph.fPath, _ = os.Getwd()
	}
	err := ph.executeSwitch()
	if err != nil {
		panic(fmt.Errorf("执行出错：%s\n", err))
	}
}
