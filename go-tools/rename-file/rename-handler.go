package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"go-tools/fileCommon"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// ============================执行重命名操作==========================================

func (f *RenameHandler) RenameOperate(oldStore, newStore string) error {
	if err := os.Rename(oldStore, newStore); err != nil {
		return err
	}
	log.Printf("重命名为 %s 成功，请查看...\n", newStore)
	return nil
}

// ============================增加文件名标志==========================================
func (f *RenameHandler) AddSign(ctx *cli.Context) error {
	for {
		select {
		case chanInfo := <-f.fileInfoChan:
			finfo := chanInfo.(map[string]interface{})
			for path, value := range finfo {
				info := value.(os.FileInfo)
				filePath := fileCommon.GetFilePath(path)
				fileExt := fileCommon.GetFileExt(info.Name())
				fileName := fileCommon.GetFileName(info.Name())
				if f.nameLoc {
					newName = fmt.Sprintf("%v-%v%s", f.addStr, fileName, fileExt)
				} else {
					newName = fmt.Sprintf("%v-%v%s", fileName, f.addStr, fileExt)
				}

				if filepath.Clean(f.outputPath) == "." {
					newStore = fmt.Sprintf("%s%s%s", filePath, pathMark, newName)
				} else {
					err, out := f.OutPutOperate(ctx)
					if err != nil {
						return err
					}
					newStore = fmt.Sprintf("%s%s%s", out, pathMark, newName)
				}
				//log.Println(newStore)
				if err := f.RenameOperate(path, newStore); err != nil {
					return err
				}
			}

		default:
			return nil

		}
	}
}

// ============================替换文件的字符串==========================================
func (f *RenameHandler) ReplaceFileName(ctx *cli.Context) error {
	for {
		select {
		case chanInfo := <-f.fileInfoChan:
			finfo := chanInfo.(map[string]interface{})
			for path, value := range finfo {
				info := value.(os.FileInfo)
				filePath := fileCommon.GetFilePath(path)
				fileName := fileCommon.GetFileName(info.Name())
				fileExt := fileCommon.GetFileExt(info.Name())
				if !strings.Contains(fileName, f.oldName) {
					continue
				}
				newName = strings.ReplaceAll(fileName, f.oldName, f.newName)
				if filepath.Clean(f.outputPath) == "." {
					newStore = fmt.Sprintf("%s%s%s%s", filePath, pathMark, newName, fileExt)
				} else {
					err, out := f.OutPutOperate(ctx)
					if err != nil {
						return err
					}
					newStore = fmt.Sprintf("%s%s%s%s", out, pathMark, newName, fileExt)
				}
				//log.Println(newStore)
				if err := f.RenameOperate(path, newStore); err != nil {
					return err
				}
			}
		default:
			return nil
		}
	}
}

// ============================使用相同名字作为文件名==========================================
func (f *RenameHandler) UseSameName(ctx *cli.Context) error {
	var (
		folderList     []string
		sameList       []string
		newList        []string
		originFileList []string
		tmpName        string
		count          = 0
	)
	for {
		select {
		case chanInfo := <-f.fileInfoChan:
			finfo := chanInfo.(map[string]interface{})
			for path, value := range finfo {
				info := value.(os.FileInfo)
				filePath := fileCommon.GetFilePath(path)
				fileExt := fileCommon.GetFileExt(info.Name())
				fileName := fileCommon.GetFileName(info.Name())
				if strings.Compare(filePath, filepath.Clean(f.inputPath)) == 1 {
					folderList = append(folderList, path)
				} else {
					count++
					if strings.Contains(fileName, "-") {
						tmpsplit := strings.Split(fileName, "-")
						tmpName = tmpsplit[0]
					}

					if f.sameFileName == tmpName {
						sameList = append(sameList, info.Name())
					} else {
						originFileList = append(originFileList, path)
					}

					newName = fmt.Sprintf("%s-%s%s", f.sameFileName, fmt.Sprintf("%02d", count), fileExt)
					newList = append(newList, newName)
					if filepath.Clean(f.outputPath) == "." {
						newStore = fmt.Sprintf("%s%s", filePath, pathMark)
					} else {
						err, out := f.OutPutOperate(ctx)
						if err != nil {
							return err
						}
						newStore = fmt.Sprintf("%s%s", out, pathMark)
					}
				}
			}
		default:
			if err := f.sameNameHandle(newList, sameList, originFileList, newStore); err != nil {
				return err
			}
			f.FolderHandle(ctx, folderList) //这里处理内层文件夹
			return nil
		}
	}
}

func (f *RenameHandler) sameNameHandle(newList, sameList, originFileList []string, tmpNewStore string) error {
	mp := make(map[string]bool)
	for _, s := range newList {
		if _, ok := mp[s]; !ok {
			mp[s] = true
		}
	}
	for _, s := range sameList {
		if _, ok := mp[s]; ok {
			delete(mp, s)
		}
	}
	var count = 0
	for newName = range mp {
		newStore = fmt.Sprintf("%s%s", tmpNewStore, newName)
		originFile := originFileList[count]
		//log.Println(originFile, newStore)
		if err := f.RenameOperate(originFile, newStore); err != nil {
			return err
		}
		count++
	}
	return nil
}

// ============================使用所在文件夹名作为新文件的名字==========================================
func (f *RenameHandler) UsePathName(ctx *cli.Context) error {
	if f.nameFormPath == false {
		return nil
	}
	count := 0
	var folderList = []string{}
	for {
		select {
		case chanInfo := <-f.fileInfoChan:
			finfo := chanInfo.(map[string]interface{})
			for path, value := range finfo {
				Info := value.(os.FileInfo)
				pathInfo := strings.Split(path, pathMark)
				pathName := strings.Join(pathInfo[len(pathInfo)-2:len(pathInfo)-1], "")
				fpath := fileCommon.GetFilePath(path)
				if strings.Compare(fpath, filepath.Clean(f.inputPath)) == 1 {
					folderList = append(folderList, path)
				} else {
					count++
					newName = fmt.Sprintf("%s-%v%s", pathName, fmt.Sprintf("%02d", count), fileCommon.GetFileExt(Info.Name()))
					if filepath.Clean(f.outputPath) == "." {
						newStore = fmt.Sprintf("%s%s%s", fpath, pathMark, newName)
					} else {
						err, out := f.OutPutOperate(ctx)
						if err != nil {
							return err
						}
						newStore = fmt.Sprintf("%s%s%s", out, pathMark, newName)
					}
					//log.Println("newstore:", newStore)
					if err := f.RenameOperate(path, newStore); err != nil {
						return err
					}
				}
			}
		default:
			f.FolderHandle(ctx, folderList)
			return nil
		}
	}
}

// ============================补齐文件前缀的数字==========================================
func (f *RenameHandler) AlterSerial(ctx *cli.Context) error {
	for {
		select {
		case chanInfo := <-f.fileInfoChan:
			fileInfo := chanInfo.(map[string]interface{})
			for path, value := range fileInfo {
				info := value.(os.FileInfo)
				filePath := fileCommon.GetFilePath(path)
				fileName := fileCommon.GetFileName(info.Name())
				fileExt := fileCommon.GetFileExt(info.Name())
				serial := f.reNum(fileName)
				sLength := len(serial)
				if sLength == 0 {
					continue
				}
				numArgs := []string{"%", "0", strconv.Itoa(f.filesn), "d"}
				digit := strings.Join(numArgs, "")
				sNumInt, _ := strconv.Atoi(serial)
				snum := fmt.Sprintf(digit, sNumInt)
				nameSplit := strings.Split(fileName, serial)
				newName = fmt.Sprintf("%s%s%s%s", nameSplit[0], snum, nameSplit[1], fileExt)
				if f.filesn == 2 && sLength == 1 {
					if filepath.Clean(f.outputPath) == "." {
						newStore = fmt.Sprintf("%s%s%s", filePath, pathMark, newName)
					} else {
						err, out := f.OutPutOperate(ctx)
						if err != nil {
							return err
						}
						newStore = fmt.Sprintf("%s%s%s", out, pathMark, newName)
					}

					//log.Println(newStore)
					if err := f.RenameOperate(path, newStore); err != nil {
						return err
					}
				}

				if f.filesn > 2 {
					if filepath.Clean(f.outputPath) == "." {
						newStore = fmt.Sprintf("%s%s%s", filePath, pathMark, newName)
					} else {
						err, out := f.OutPutOperate(ctx)
						if err != nil {
							return err
						}
						newStore = fmt.Sprintf("%s%s%s", out, pathMark, newName)
					}
					//log.Println(newStore)
					if err := f.RenameOperate(path, newStore); err != nil {
						return err
					}
				}

			}
		default:
			return nil
		}
	}
}

func (f *RenameHandler) reNum(title string) string {
	reg, _ := regexp.Compile("\\d+")
	return reg.FindString(title)
}

// ============================删除文件名中的某个字符==========================================
func (f *RenameHandler) SubFileName(ctx *cli.Context) error {
	for {
		select {
		case chanInfo := <-f.fileInfoChan:
			fileInfo := chanInfo.(map[string]interface{})
			for path, value := range fileInfo {
				info := value.(os.FileInfo)
				fileName := fileCommon.GetFileName(info.Name())
				fileExt := fileCommon.GetFileExt(info.Name())
				if f.subStr == "" || !strings.Contains(fileName, f.subStr) {
					continue
				}
				fpath := fileCommon.GetFilePath(path)
				tmpName := []rune(fileName)
				switch {
				case f.subLoc == "all":
					newName = strings.ReplaceAll(fileName, f.subStr, "")
					if filepath.Clean(f.outputPath) == "." {
						newStore = fmt.Sprintf("%s%s%s%s", fpath, pathMark, newName, fileExt)
					} else {
						err, out := f.OutPutOperate(ctx)
						if err != nil {
							return err
						}
						newStore = fmt.Sprintf("%s%s%s%s", out, pathMark, newName, fileExt)
					}
					//log.Println(newStore)
					if err := f.RenameOperate(path, newStore); err != nil {
						return err
					}
				case f.subLoc == "left":
					if string(tmpName[0]) == f.subStr {
						newName = strings.TrimLeft(fileName, f.subStr)
						if filepath.Clean(f.outputPath) == "." {
							newStore = fmt.Sprintf("%s%s%s%s", fpath, pathMark, newName, fileExt)
						} else {
							err, out := f.OutPutOperate(ctx)
							if err != nil {
								return err
							}
							newStore = fmt.Sprintf("%s%s%s%s", out, pathMark, newName, fileExt)
						}

						//log.Println(newStore)
						if err := f.RenameOperate(path, newStore); err != nil {
							return err
						}
					}

				case f.subLoc == "right":
					if string(tmpName[len(tmpName)-1]) == f.subStr {
						newName = strings.TrimRight(fileName, f.subStr)
						if filepath.Clean(f.outputPath) == "." {
							newStore = fmt.Sprintf("%s%s%s%s", fpath, pathMark, newName, fileExt)
						} else {
							err, out := f.OutPutOperate(ctx)
							if err != nil {
								return err
							}
							newStore = fmt.Sprintf("%s%s%s%s", out, pathMark, newName, fileExt)
						}
						//log.Println(newStore)
						if err := f.RenameOperate(path, newStore); err != nil {
							return err
						}
					}

				default:
					return nil
				}
			}
		default:
			return nil
		}
	}
}

// ============================删除文件名中的空格==========================================
func (f *RenameHandler) RmSpace(ctx *cli.Context) error {
	for {
		select {
		case chanInfo := <-f.fileInfoChan:
			finfo := chanInfo.(map[string]interface{})
			for path, value := range finfo {
				info := value.(os.FileInfo)
				nameInfo := fileCommon.GetFileName(info.Name())
				if strings.Contains(nameInfo, " ") {
					newName = strings.Replace(nameInfo, " ", "", -1)
				} else {
					newName = nameInfo
				}

				if !strings.Contains(nameInfo, " ") {
					continue
				}
				if filepath.Clean(f.outputPath) == "." {
					newStore = fmt.Sprintf("%s%s%s%s", fileCommon.GetFilePath(path), pathMark, newName, fileCommon.GetFileExt(info.Name()))

				} else {
					err, out := f.OutPutOperate(ctx)
					if err != nil {
						return err
					}
					newStore = fmt.Sprintf("%s%s%s%s", out, pathMark, newName, fileCommon.GetFileExt(info.Name()))
				}
				//log.Println(newStore)
				if err := f.RenameOperate(path, newStore); err != nil {
					return err
				}
			}
		default:
			return nil
		}
	}
}
