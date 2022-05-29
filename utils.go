package main

import (
	"errors"
	"os"
	"strings"
	"time"
)

const (
	DirPerm = 0755
)

func IsFolderExist(path string) (isExist bool, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// folder not exist
			return false, nil
		}

		// other error
		return false, err
	}

	// check if it is a folder
	if !stat.IsDir() {
		err = errors.New("it's not folder")
		return false, err
	}

	// folder exists
	return true, nil
}
func CreateFolder(path string) error {
	err := os.MkdirAll(path, DirPerm)
	if err != nil {
		return err
	}
	return nil
}

func GetTimestamp() string {
	return strings.ReplaceAll(time.Now().Format("2006_01_02__15_04_05.999999999"), ".", "__")
}
