package fileutils

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

func CreateDirIfNotExists(path string) {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path, 0775)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func CreateFileIfNotExists(path string) *os.File {
	dir := filepath.Dir(path)
	CreateDirIfNotExists(dir)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		log.Fatalln(err)
	}
	return file
}
