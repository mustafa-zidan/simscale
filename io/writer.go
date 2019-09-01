package io

import (
	"bufio"
	"encoding/json"
	"fmt"
	. "github.com/mustafa-zidan/simscale/domain"
	"log"
	"os"
)

type Writer interface {
	Append(*LogTree)
	Open(string)
	Close()
}

type FileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

func (f *FileWriter) Open(filePath string) *FileWriter {
	if ok, _ := f.Exists(filePath); !ok {
		_, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("failed creating to file: %s", err)
		}
	}
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	f.file = file
	f.writer = bufio.NewWriter(file)
	return f
}

func (f *FileWriter) Close() {
	f.writer.Flush()
	f.file.Close()

}

func (f *FileWriter) Append(logtree *LogTree) {
	lt, err := json.Marshal(logtree)
	if err != nil {
		log.Fatalf("failed marshalling logtree to JSON: %s", err)
	}
	f.writer.WriteString(fmt.Sprintf("%s\n", lt))
	f.writer.Flush()
}

func (f *FileWriter) Exists(file string) (bool, error) {
	_, err := os.Stat(file)
	return !os.IsNotExist(err), err
}

func NewFileWriter(filePath string) *FileWriter {
	f := &FileWriter{}
	f.Open(filePath)
	return f
}
