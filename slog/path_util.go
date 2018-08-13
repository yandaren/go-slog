// File path_util.go
// @Author: yandaren1220@126.com
// @Date: 2018-08-13

package slog

import (
	"os"
)

// check path exist
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// remove file)
func RemoveFile(path string) (bool, error) {
	err := os.Remove(path)
	if err != nil {
		return false, err
	}
	return true, nil
}

// remove path
func RemovePath(path string) (bool, error) {
	err := os.RemoveAll(path)
	if err != nil {
		return false, err
	}
	return true, nil
}

// rename path
func RenamePath(src string, target string) (bool, error) {
	err := os.Rename(src, target)
	if err != nil {
		return false, err
	}
	return true, nil
}
