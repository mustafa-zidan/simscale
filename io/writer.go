package io

import (
	"encoding/json"
	"log"
	"os"
)

type Writer interface {
	Append(*LogTree)
	Open(string)
	Close()
}

type FileWriter struct {
	file *os.File
}

func (f *FileWriter) Open(filePath string) *FileWriter {
	f = &FileWriter{}
	if ok, _ := f.Exists(filePath); !ok {
		_, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("failed creating to file: %s", err)
		}
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	f.file = file
	return f
}

func (f *FileWriter) Close() {
	f.file.Close()
}

func (f *FileWriter) Append(logtree *LogTree) {
	e, err := json.Marshal(logtree)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
	_, err = f.file.WriteString(string(e))
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}

func (f *FileWriter) Exists(file string) (bool, error) {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}
