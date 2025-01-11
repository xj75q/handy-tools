package colorLog

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"time"
)

func NewLog(opts ...logOption) *logParam {
	logOpt := &logParam{
		level:    debugLevel,
		saveDays: 30,
		isShort:  false,
		saveMode: byDay,
		fileName: defaultName,
		filePath: currentPath,
	}
	for _, opt := range opts {
		opt(logOpt)
	}
	log.SetOutput(logOpt)
	return logOpt
}

func (f *logParam) withColor(text string, colour int64) string {
	col := color.New(colors[colour]...)
	return col.Sprint(text)
}

func (f *logParam) Write(buf []byte) (n int, err error) {
	store := newLogStore()
	switch f.saveMode {
	case bySize:
		fileInfo, err := os.Stat(f.fileName)
		if err != nil {
			store.createLogFile(f)
			f.logTime = time.Now().Unix()
		} else {
			filesize := fileInfo.Size()
			if f.filePtr == nil || filesize > f.filesize {
				store.createLogFile(f)
				f.logTime = time.Now().Unix()
			}
		}
	default: // 默认按天  ByDay
		if f.logTime+3600 < time.Now().Unix() {
			if err = store.createLogFile(f); err != nil {
				return 0, err
			}
			f.logTime = time.Now().Unix()
		}
	}

	if f.filePtr == nil {
		fmt.Printf("log filePtr is nil !\n")
		return len(buf), nil
	}
	return f.filePtr.Write(buf)
}

// ======================实现的接口方法===========================================

func (f *logParam) Debugf(format string, args ...interface{}) {
	log.SetPrefix(f.withColor(debugPre, debugLevel))
	log.Output(2, fmt.Sprintf("%s\n", fmt.Sprintf(format, args...)))
}

func (f *logParam) Infof(format string, args ...interface{}) {
	log.SetPrefix(f.withColor(infoPre, infoLevel))
	log.Output(2, fmt.Sprintf("%s\n", fmt.Sprintf(format, args...)))
}

func (f *logParam) Warnf(format string, args ...interface{}) {
	log.SetPrefix(f.withColor(warnPre, warnLevel))
	log.Output(2, fmt.Sprintf("%s\n", fmt.Sprintf(format, args...)))
}

func (f *logParam) Errorf(format string, args ...interface{}) {
	log.SetPrefix(f.withColor(errorPre, errorLevel))
	log.Output(2, fmt.Sprintf("%s\n", fmt.Sprintf(format, args...)))
}

func (f *logParam) Fatalf(format string, args ...interface{}) {
	log.SetPrefix(f.withColor(fatalPre, fatalLevel))
	log.Output(2, fmt.Sprintf("%s\n", fmt.Sprintf(format, args...)))
}
