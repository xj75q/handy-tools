package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"go-tools/fileCommon"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

func (f *RenameHandler) FolderHandle(ctx *cli.Context, originPathList []string) {
	keys := []string{}
	for _, originPath := range originPathList {
		folderPath := fileCommon.GetFilePath(originPath)
		keys = append(keys, folderPath)
	}

	sort.Strings(keys)
	result := make(map[string]int)
	for _, item := range keys {
		result[item]++
	}

	var numList = []map[string]int{}
	for key, value := range result {
		var vList = make([]int, value)
		for index := range vList {
			var fileNum = make(map[string]int)
			fileNum[key] = index + 1
			numList = append(numList, fileNum)
		}
	}
	go f.generateFileNum(ctx, originPathList, numList)
	f.manageLayer(originPathList)
}

func (f *RenameHandler) generateFileNum(ctx *cli.Context, originPathList []string, numList []map[string]int) {
	for _, originPath := range originPathList {
		pathInfo := strings.Split(originPath, pathMark)
		fname := strings.Join(pathInfo[len(pathInfo)-2:len(pathInfo)-1], pathMark)
		fileExt := fileCommon.GetFileExt(originPath)
		filePath := fileCommon.GetFilePath(originPath)

		for _, fileNum := range numList {
			if fileNum[filePath] != 0 {
				//todo 这里更改为适合各种参数
				if f.nameFormPath == true {
					newName = fmt.Sprintf("%s-%v%s", fname, fmt.Sprintf("%02d", fileNum[filePath]), fileExt)
					if filepath.Clean(f.outputPath) == "." {
						newStore = fmt.Sprintf("%s%s%s", filePath, pathMark, newName)
						f.newFileChan <- newStore
					} else {
						err, out := f.OutPutOperate(ctx)
						if err != nil {
							//return err
							return
						}
						newStore = fmt.Sprintf("%s%s%s", out, pathMark, newName)
						f.newFileChan <- newStore
					}

				} else {
					if filepath.Clean(f.outputPath) == "." {
						newStore = fmt.Sprintf("%s%s%s-%v%s", filePath, pathMark, f.sameFileName, fmt.Sprintf("%02d", fileNum[filePath]), fileExt)
						f.newFileChan <- newStore
					} else {
						err, out := f.OutPutOperate(ctx)
						if err != nil {
							//return err
						}
						newStore = fmt.Sprintf("%s%s%s-%v%s", out, pathMark, f.sameFileName, fmt.Sprintf("%02d", fileNum[filePath]), fileExt)
						f.newFileChan <- newStore
					}
					newStore = fmt.Sprintf("%s%s%s-%v%s", filePath, pathMark, f.sameFileName, fmt.Sprintf("%02d", fileNum[filePath]), fileExt)
					f.newFileChan <- newStore
				}
				delete(fileNum, filePath)
			}
		}
	}

}

func (f *RenameHandler) manageLayer(originPathList []string) {
	for _, originPath := range originPathList {
		select {
		case newStore = <-f.newFileChan:
			//log.Println(newStore)
			//fmt.Sprintln(originPath)
			if err := f.RenameOperate(originPath, newStore); err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

func (f *RenameHandler) SameNameHandle() {

}
