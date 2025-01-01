package colorLog

import (
	"github.com/fatih/color"
	"os"
	"time"
)

const (
	//定义日志前缀
	debugPre = "[DEBUG]"
	infoPre  = "[INFO] "
	warnPre  = "[WARN] "
	errorPre = "[ERROR]"
	fatalPre = "[FATAL]"

	//定义日志等级
	debugLevel int64 = iota + 1
	infoLevel
	warnLevel
	errorLevel
	fatalLevel

	//定义日志
	byDay int64 = iota + 1
	bySize
)

var (
	newFile        string
	logFile        string
	pathFlag       = string(os.PathSeparator)
	defaultName    = time.Now().Format("20060102")
	currentPath, _ = os.Getwd()
	colors         = map[int64][]color.Attribute{
		debugLevel: {color.FgHiWhite},
		infoLevel:  {color.FgYellow},
		warnLevel:  {color.FgBlue},
		errorLevel: {color.FgRed},
		fatalLevel: {color.FgGreen},
	}
)

type (
	//如果需要对外提供，将小写改为大写
	colorLog interface {
		Debugf(format string, args ...interface{})
		Infof(format string, args ...interface{})
		Warnf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
		Fatalf(format string, args ...interface{})
	}

	logStore struct{}

	logOption func(logopt *logParam)

	logParam struct {
		filePath string
		fileName string

		level    int64
		saveMode int64 // 保存模式
		filesize int64 // 设置 saveMode 为 bySize 生效
		saveDays int64 // 最大存活天数
		logTime  int64 //日志里记录的时间
		isShort  bool  //使用短名字

		filePtr *os.File
	}
)
