package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	noteExecPath, _ = os.Executable()
	notePath        = filepath.Dir(noteExecPath)
	notePathFlag    = string(os.PathSeparator)
	notecfgName     = "notecfg.json"
)

type notescfg struct {
	Path     string `json:"path"`
	Filename string `json:"filename"`
	FileType bool   `json:"filetype"`
}

func noteHandler() *notescfg {
	return &notescfg{}
}

func (c *notescfg) initConfig() (error, string) {
	file := filepath.Clean(fmt.Sprintf("%s%s%s", notePath, notePathFlag, notecfgName))
	_, err := os.Stat(file)
	if err != nil && os.IsNotExist(err) {
		createFile, _ := os.Create(file)
		rb, _ := json.Marshal(c)
		_, err := createFile.Write(rb)
		if err != nil {
			return fmt.Errorf("创建并写入文件失败，请检查..."), ""
		} else {
			return nil, fmt.Sprintf(">> 文本记录配置文件初始化成功...")
		}
	} else {
		rb, _ := json.Marshal(c)
		if err := ioutil.WriteFile(file, rb, 0644); err != nil {
			return fmt.Errorf(">> 更新配置文件失败，请重试..."), ""
		} else {
			return nil, fmt.Sprintf(">> 更新配置文件成功..")
		}
	}
	return nil, ""
}

func (c *notescfg) readCfg() (error, *notescfg) {
	file := filepath.Clean(fmt.Sprintf("%s%s%s", notePath, notePathFlag, notecfgName))
	_, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("配置文件不存在"), nil
	}

	cfgName := strings.Split(notecfgName, ".")
	viper.SetConfigName(cfgName[0])
	viper.SetConfigType(cfgName[1])
	viper.AddConfigPath(notePath)
	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("读取配置文件出错：%v\n", err), nil
	}
	if err = viper.Unmarshal(&c); err != nil {
		return fmt.Errorf("文件解析出错：%v\n", err), nil
	}
	return nil, c
}

func (n *notescfg) recordNote(content string) error {
	var notefile string
	if n.FileType == true {
		notefile = filepath.Clean(fmt.Sprintf("%s%s%s.txt", n.Path, notePathFlag, n.Filename))
	} else {
		notefile = filepath.Clean(fmt.Sprintf("%s%s%s", n.Path, notePathFlag, n.Filename))
	}
	file, err := os.OpenFile(notefile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("打开文件失败:", err)
	}
	defer file.Close()
	insertData := fmt.Sprintf("\n\n%v", content)
	_, err = file.Write([]byte(insertData))
	if err != nil {
		return fmt.Errorf("写入文件失败:", err)
	}
	return nil
}

var (
	note    = noteHandler()
	noteCmd = &cobra.Command{
		Short: "随手记录文本信息",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	cfgCmd = &cobra.Command{
		Use:     "config",
		Short:   "设置配置信息",
		Aliases: []string{"cfg"},
		Run: func(cmd *cobra.Command, args []string) {
			path, _ := cmd.Flags().GetString("path")
			filename, _ := cmd.Flags().GetString("filename")
			fileType, _ := cmd.Flags().GetBool("type")
			if path == "" || filename == "" {
				log.Println(">> 请填写正确的配置信息")
				return
			}

			if !strings.Contains(path, notePathFlag) {
				log.Println("配置信息路径不正确，请重新填写...")
				return
			}

			note.Path = path
			note.Filename = filename
			note.FileType = fileType
			err, result := note.initConfig()
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(result)
		},
	}

	recordCmd = &cobra.Command{
		Use:   "note",
		Short: "记录文本",
		RunE: func(cmd *cobra.Command, args []string) error {
			content, _ := cmd.Flags().GetString("content")
			err, cfg := note.readCfg()
			if err != nil {
				return err
			}
			note.Path = cfg.Path
			note.Filename = cfg.Filename
			if err := note.recordNote(content); err != nil {
				return err
			}
			return nil
		},
	}
)

func init() {
	cfgCmd.Flags().StringP("path", "p", "", "填写笔记保存的路径")
	cfgCmd.Flags().StringP("filename", "n", "", "填写笔记名称")
	cfgCmd.Flags().BoolP("type", "t", true, "笔记是否需要后缀.txt")
	recordCmd.Flags().StringP("content", "c", "", "需记录的内容")

	noteCmd.AddCommand(cfgCmd)
	noteCmd.AddCommand(recordCmd)
}

func main() {
	if err := noteCmd.Execute(); err != nil {
		log.Printf(">> 出错了：%s\n", err)
		os.Exit(1)
	}
}
