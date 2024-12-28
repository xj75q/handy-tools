package fileCommon

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

var pathMark = string(os.PathSeparator)

func IsFile(path string) bool {
	return !checkDir(path)
}

func checkDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func GetFileExt(param string) string {
	return path.Ext(param)
}

func GetFileName(param string) (fileName string) {
	ext := GetFileExt(param)
	if strings.Contains(param, pathMark) {
		tmpFile := strings.Split(param, pathMark)
		fileName = strings.TrimSuffix(tmpFile[len(tmpFile)-1:][0], ext)
		return fileName
	} else {
		fileName = strings.TrimSuffix(param, ext)
		return fileName
	}
}

func GetFilePath(param string) string {
	tmpFile := strings.Split(filepath.Clean(param), pathMark)
	fpath := strings.Join(tmpFile[:len(tmpFile)-1], pathMark)
	return fpath
}

func CheckSavePath(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsNotExist(err)
}

func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsPermission(err)
}

func CreateSavePath(dst string, perm os.FileMode) error {
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}
	return nil
}
