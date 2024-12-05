package fileCommon

import "os"

func IsFile(path string) bool {
	return !isDir(path)
}

func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
