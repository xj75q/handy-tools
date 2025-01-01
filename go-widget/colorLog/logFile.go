package colorLog

import (
	"fmt"
	"go-widget/fileCommon"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ======================Set设置===========================================

func SetLog(opts ...logOption) logOption {
	return func(logOpt *logParam) {
		for _, opt := range opts {
			opt(logOpt)
		}
	}
}

func SetPath(path string) logOption {
	return func(opt *logParam) {
		if path == "" {
			opt.filePath = currentPath
			return
		}
		opt.filePath = path
	}
}

func SetfileName(name string) logOption {
	return func(opt *logParam) {
		if name == "" {
			opt.fileName = defaultName
			return
		}
		opt.fileName = name
	}
}

func SetLevel(level int64) logOption {
	return func(opt *logParam) {
		opt.level = level
	}
}

// size或day模式
func SetsaveMode(saveMode int64) logOption {
	return func(opt *logParam) {
		opt.saveMode = saveMode
	}
}

func SetInfoLoc(isTrue bool) logOption {
	return func(opt *logParam) {
		if isTrue {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
		} else {
			log.SetFlags(log.LstdFlags)
		}
	}
}

func SetMaxSize(size int64) logOption {
	return func(opt *logParam) {
		opt.filesize = size
	}
}

func SetMaxAge(day int64) logOption {
	return func(opt *logParam) {
		opt.saveDays = day
	}
}

// =====================处理logFile文件===========================================
func newLogStore() *logStore {
	return &logStore{}
}

func (s *logStore) isToday(d time.Time) bool {
	now := time.Now()
	return d.Year() == now.Year() && d.Month() == now.Month() && d.Day() == now.Day()
}

func (s *logStore) createLogFile(f *logParam) (err error) {
	err, logFile = s.handleFile(f)
	if err != nil {
		return err
	}
	/*todo 这里需优化，
	执行程序时检测路径下的文件，
	将超过最大时间后进行压缩并放入新文件夹
	在bysize的时候，旧文件重命名并进行新文件创建
	*/
	now := time.Now()
	originFile := fmt.Sprintf("%s%s%s", filepath.Clean(f.filePath), pathFlag, f.fileName)
	newName := fmt.Sprintf("%s_%02d%02d_old.log", originFile, now.Month(), now.Day())
	_, err = os.Stat(logFile)
	if !(err != nil && os.IsNotExist(err)) && !s.isToday(now) {
		if err = os.Rename(logFile, newName); err != nil {
			return fmt.Errorf("重命名log文件失败：%v", err)
		}
	}

	s.openfile(f)
	return nil
}

func (s *logStore) handleFile(f *logParam) (error, string) {
	name := f.fileName
	cleanPath := filepath.Clean(f.filePath)
	if fileCommon.CheckSavePath(cleanPath) {
		if err := fileCommon.CreateSavePath(cleanPath, os.FileMode(0755)); err != nil {
			return fmt.Errorf("创建目录失败![%v]\n", err), ""
		}
	}
	permission := fileCommon.CheckPermission(cleanPath)
	if !permission {
		if err := os.Chmod(cleanPath, os.FileMode(0755)); err != nil {
			return fmt.Errorf("目录赋权限失败![%v]\n", err), ""
		}
	}
	chars := `\ / ： * ？" < > |`
	if strings.ContainsAny(name, chars) {
		return fmt.Errorf("文件名含有无效字符，请重新命名！！！"), ""
	}
	logFile = fmt.Sprintf("%s%s%s.log", cleanPath, pathFlag, name)
	if _, err := os.Stat(logFile); err != nil {
		if os.IsNotExist(err) {
			if err = fileCommon.CreateFile(logFile); err != nil {
				return err, ""
			}
		}
	}
	return nil, logFile
}

func (s *logStore) openfile(l *logParam) {
	for i := 0; i < 10; i++ {
		if fd, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); nil == err {
			l.filePtr.Sync()
			l.filePtr.Close()
			l.filePtr = fd
			break
		} else {
			fmt.Printf("打开log文件出错：%v\n ", err)
		}
		l.filePtr = nil
	}
}
