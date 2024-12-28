package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"go-tools/fileCommon"
	"os"
	"path/filepath"
)

var (
	newName  string
	newStore string
	pathMark = string(os.PathSeparator)
)

type (
	RenameHandler struct {
		SameField
		ReplaceField
		SubField
		AltersnField
		AddsignField
		BoolField

		inputPath  string
		outputPath string

		fileInfoChan chan interface{}
		folderChan   chan interface{}
		newFileChan  chan string
	}

	SameField struct {
		sameFileName string
	}

	ReplaceField struct {
		oldName string
		newName string
	}

	SubField struct {
		subStr string
		subLoc string
	}

	AltersnField struct {
		filesn int
	}

	AddsignField struct {
		addStr  string
		nameLoc bool
	}

	BoolField struct {
		nameFormPath bool
		removeSpace  bool
	}
)

func NewRename() *RenameHandler {
	return &RenameHandler{
		fileInfoChan: make(chan interface{}, 20),
		folderChan:   make(chan interface{}, 20),
		newFileChan:  make(chan string, 20),
	}
}

func (f *RenameHandler) InputOperate(ctx *cli.Context) error {
	err := filepath.Walk(f.inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fpath := filepath.Clean(fmt.Sprintf("%s%s%s", f.inputPath, pathMark, info.Name()))
		outputpath := filepath.Clean(f.outputPath)
		if info.IsDir() && fpath == outputpath {
			return nil
		}
		if fileCommon.IsFile(path) {
			fileInfo := make(map[string]interface{})
			fileInfo[path] = info
			go func(ctx *cli.Context) {
				f.fileInfoChan <- fileInfo
			}(ctx)
		}
		return nil
	})
	return err
}

func (f *RenameHandler) OutPutOperate(ctx *cli.Context) (error, string) {
	out := filepath.Clean(f.outputPath)
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
