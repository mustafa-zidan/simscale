package io

import (
	"bufio"
	"log"
	"os"
)

type Reader interface {
	Open()
	Read(chan string)
	Close()
}

type FileReader struct {
	file *os.File
}

func (f *FileReader) Open(filePath string) *FileReader {
	if ok, _ := f.Exists(filePath); !ok {
		log.Fatalf("file not found create empty: %s", filePath)
		_, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("failed creating to file: %s", err)
		}
	}
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	f.file = file
	return f
}

func (f *FileReader) Close() {
	f.file.Close()
}

func (f *FileReader) Read(lines chan string) {
	scanner := bufio.NewScanner(f.file)
	for scanner.Scan() {
		// Later I want to create a buffer of lines,
		// not just line-by-line here ...
		lines <- scanner.Text()
	}
	close(lines)
}

func (f *FileReader) Exists(file string) (bool, error) {
	_, err := os.Stat(file)
	return !os.IsNotExist(err), err
}

func NewFileReader(filePath string) *FileReader {
	f := &FileReader{}
	f.Open(filePath)
	return f
}
