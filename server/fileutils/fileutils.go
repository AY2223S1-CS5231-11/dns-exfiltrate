package fileutils

import (
	"errors"
	"log"
	"os"
)

func CreateDirIfNotExists(path string) {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(path, 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func CreateFileIfNotExists(path string) *os.File {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	return file
}
