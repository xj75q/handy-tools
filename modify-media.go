package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	batchsize                   = 5
	channelsize                 = 10
	timeout       time.Duration = 1
	commandFfmpeg               = "ffmpeg"
	pathFlag                    = string(os.PathSeparator)
)

type (
	batcher struct {
		option     batchOption
		wg         *sync.WaitGroup
		quit       chan struct{}
		eventQueue chan interface{}
	}

	batchOption struct {
		BatchSize   int
		Workers     int
		ChannelSize int
		DelayTime   time.Duration
	}
)

func batchHandler() *batcher {
	option := batchOption{
		BatchSize:   batchsize,
		Workers:     runtime.NumCPU(),
		ChannelSize: channelsize,
		DelayTime:   timeout * time.Second,
	}

	if option.ChannelSize <= option.Workers {
		option.ChannelSize = option.Workers
	}

	return &batcher{
		option:     option,
		wg:         new(sync.WaitGroup),
		quit:       make(chan struct{}),
		eventQueue: make(chan interface{}, option.ChannelSize),
	}
}

func (b *batcher) execute() {
	b.wg.Add(1)
	defer b.wg.Done()
	batch := make([]interface{}, 0, b.option.BatchSize)
	delayTimer := time.NewTimer(0)
	if !delayTimer.Stop() {
		<-delayTimer.C
	}
	defer delayTimer.Stop()
LOOP:
	for {
		select {
		case req := <-b.eventQueue:
			batch = append(batch, req)
			bsize := len(batch)
			if bsize < b.option.BatchSize {
				if bsize == 1 {
					delayTimer.Reset(b.option.DelayTime)
				}
				break
			}

			b.batchProcess(batch)
			if !delayTimer.Stop() {
				<-delayTimer.C
			}
			batch = make([]interface{}, 0, b.option.BatchSize)

		case <-delayTimer.C:
			if len(batch) == 0 {
				break
			}
			b.batchProcess(batch)
			batch = make([]interface{}, 0, b.option.BatchSize)

		case <-b.quit:
			if len(batch) > 0 {
				b.batchProcess(batch)
			}
			batch = make([]interface{}, 0, b.option.BatchSize)
			break LOOP

		default:

		}

	}

}

func (b *batcher) stop() {
	if len(b.eventQueue) == 0 {
		close(b.quit)
	} else {
		ticker := time.NewTicker(50 * time.Millisecond)
		for range ticker.C {
			if len(b.eventQueue) == 0 {
				close(b.quit)
				break
			}
		}
		ticker.Stop()
	}
	b.wg.Wait()
}

func (b *batcher) batchProcess(items []interface{}) {
	finfo := inputHandler()
	length := len(items)
	if length == 0 {
		return
	}
	for i := 0; i < length; i++ {
		b.wg.Add(1)
		go func(item interface{}) {
			defer b.wg.Done()
			finfo.switchFile(item)
		}(items[i])
	}
}

func (b *batcher) StartSwitch() {
	now := time.Now()
	fmt.Println(">> 正在转换，请稍后...")
	defer func() {
		cost := time.Since(now).String()
		fmt.Printf("总耗时为：%s\n", cost)
	}()
	defer b.stop()
	for i := 1; i <= b.option.Workers; i++ {
		go b.execute()
	}
}

type param struct {
	filePath string
	outPut   string
	ftype    string
	speed    float64
	volume   int64
}

func inputHandler() *param {
	return &param{}
}

func (p *param) isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func (p *param) isFile(path string) bool {
	return !p.isDir(path)
}

func (p *param) judgeVideoType(flag string) bool {
	switch flag {
	case "mp4":
		return true
	case "wmv":
		return true
	case "avi":
		return true
	case "rmvb":
		return true
	default:
		return false
	}
}

func (p *param) judgeAudioType(flag string) bool {
	switch flag {
	case "mp3":
		return true
	case "wav":
		return true
	case "amr":
		return true
	case "3gp":
		return true
	default:
		return false
	}
}

func (p *param) switchVideo() error {
	if p.speed != 0 && p.ftype == "video" {
		fmt.Println(">> 速度调节只可用于音频文件，请输入音频类型 audio ,之后重新运行...")
		os.Exit(0)
	}
	batch := batchHandler()
	err := filepath.Walk(p.filePath, func(pathAndFilename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		infoName := strings.Split(strings.ToLower(info.Name()), ".")
		flag := infoName[len(infoName)-1]
		switch p.ftype {
		case "video":
			isVideoType := p.judgeVideoType(flag)
			if p.isFile(pathAndFilename) && isVideoType {
				fInfo := make(map[string]interface{})
				fInfo[pathAndFilename] = p
				go func() {
					batch.eventQueue <- fInfo
				}()

				return nil
			}
		case "audio":
			isAuidoType := p.judgeAudioType(flag)
			if p.isFile(pathAndFilename) && isAuidoType {
				fInfo := make(map[string]interface{})
				fInfo[pathAndFilename] = p
				go func() {
					batch.eventQueue <- fInfo
				}()
				return nil
			}

		}
		return nil
	})
	batch.StartSwitch()
	return err
}

func (p *param) generateCmdStr(input *param, inName, outName string) (args []string) {
	if input.ftype == "video" {
		if input.volume == 0 {
			args = []string{"-i", inName, outName}
			return
		} else {
			args = []string{"-i", inName, "-filter:a", fmt.Sprintf("volume=%vdB", input.volume), outName}
			return
		}

	} else if input.ftype == "audio" {
		if input.volume != 0 && input.speed == 0 {
			args = []string{"-i", inName, "-filter:a", fmt.Sprintf("volume=%vdB", input.volume), outName}
			return
		} else if input.volume == 0 && input.speed != 0 {
			args = []string{"-i", inName, "-filter:a", fmt.Sprintf("atempo=%v", input.speed), "-vn", outName}
			return
		}
	}
	return
}

func (p *param) switchFile(fInfo interface{}) {
	finfo := fInfo.(map[string]interface{})
	for key, value := range finfo {
		info := value.(*param)
		outlist := strings.Split(key, pathFlag)
		pathInfo := strings.Join(outlist[:len(outlist)-1], pathFlag)
		fin := strings.Join(outlist[len(outlist)-1:], "")
		name := strings.Split(fin, ".")
		var (
			outName string
			fName   string
		)
		if info.ftype == "video" {
			fName = strings.Join(name[:len(name)-1], "") + ".mp3"
		} else {
			fName = "alter-" + strings.Join(name[:len(name)-1], "") + ".mp3"
		}

		if len(info.outPut) != 0 {
			outName = info.outPut + pathFlag + fName
		} else {
			outName = pathInfo + pathFlag + fName
		}

		cmdArgs := p.generateCmdStr(info, key, outName)

		cmd := exec.Command(commandFfmpeg, cmdArgs...)
		stdout, err := cmd.StdoutPipe()
		cmd.Stderr = cmd.Stdout
		if err != nil {
			fmt.Println(err)
		}
		if err = cmd.Start(); err != nil {
			fmt.Println(err)
		}

		for {
			tmp := make([]byte, 1024)
			_, err := stdout.Read(tmp)
			data := string(tmp)

			if strings.Contains(data, "muxing overhead") {
				fmt.Printf(">> 将 %s 目录下的文件转换成 %s 成功...\n", pathInfo, fName)
			}
			if err != nil {
				break
			}
		}

		if err := cmd.Wait(); err != nil {
			fmt.Println(err.Error())
		}

	}

}

var (
	rootCmd = &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			input := inputHandler()
			inpath, _ := cmd.Flags().GetString("input")
			if inpath == "./" {
				input.filePath, _ = os.Getwd()
			} else {
				input.filePath = inpath
			}

			input.outPut, _ = cmd.Flags().GetString("output")
			input.ftype, _ = cmd.Flags().GetString("type")
			input.volume, _ = cmd.Flags().GetInt64("volume")
			input.speed, _ = cmd.Flags().GetFloat64("speed")
			input.switchVideo()
		},
	}
)

func init() {
	rootCmd.Flags().StringP("input", "i", "", "文件路径(必填)")
	rootCmd.Flags().StringP("output", "o", "", "文件保存路径（选填）")
	rootCmd.Flags().StringP("type", "t", "video", "默认视频，音频请使用audio参数")
	rootCmd.Flags().Int64P("volume", "v", 0, "声音调节,声音调整分[正负数调整分贝]")
	rootCmd.Flags().Float64P("speed", "s", 0, "速度调节，倍率调整范围为[0.5, 2.0]")
	rootCmd.MarkFlagRequired("input")

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("err of %s", err)
		os.Exit(1)
	}
}
