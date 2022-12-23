package controller

import (
	"io/fs"
	"os"
)

func createFile(path string, data []byte, perm int) error {
	err := os.WriteFile(path, data, fs.FileMode(perm))
	if err != nil {
		return err
	}

	err = os.Chown(path, os.Getuid(), os.Getgid())
	return err
}
