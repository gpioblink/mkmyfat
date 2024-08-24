package tools

import (
	"fmt"
	"os"
)

func CreateSpecificatedSizeFileWhenNotExsisted(imgPath string, size uint64) (*os.File, error) {
	if Exists(imgPath) {
		return nil, fmt.Errorf("imgPath %s is exsisted. Please remove first", imgPath)
	}

	f, err := os.Create(imgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %s", imgPath, err)
	}
	defer f.Close()

	// 容量の確保
	err = os.Truncate(imgPath, int64(size))
	if err != nil {
		return f, fmt.Errorf("failed to truncate file %s: %s", imgPath, err)
	}
	return f, nil
}

func GetFileSize(filename string) (uint64, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to get file size: %s", err)
	}
	return uint64(fi.Size()), nil
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
