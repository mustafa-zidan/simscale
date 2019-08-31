package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

var (
	pattern  = "(?P<start>.*)\\s(?P<end>.*)\\s(?P<trace>.*)\\s(?P<service>.*)\\s(?P<caller_span>.*)->(?P<span>.*)"
	re       = regexp.MustCompile(pattern)
	expNames = re.SubexpNames()
)

// Parser is the interface for raw logs
type Parser interface {
	Process() error
	Parse(string) (*Log, error)
}

type RawLogParser struct {
	inputFile *os.File
}

func (p *RawLogParser) Process() error {
	defer p.close()

	scanner := bufio.NewScanner(p.inputFile)
	for scanner.Scan() {
		//TODO: Add counter lines processed
		Increment("total", 1)
		l, err := p.Parse(scanner.Text()) // the line
		if err != nil {
			log.Println(err)
		}
		InsertLogToCache(l)
	}
	log.Println(CounterList())
	return nil
}

// Parse takes a line then parse it using regex to Log Object
func (p *RawLogParser) Parse(line string) (*Log, error) {
	matches := re.FindAllStringSubmatch(line, -1)
	if len(matches) == 0 {
		Increment("ignored", 1)
		return nil, fmt.Errorf("Invalid Log Format %s, Skipping!", line)
	}

	m := map[string]string{}
	for i, n := range matches[0] {
		m[expNames[i]] = n
	}
	return NewLog(m)
}

func NewParser(path string) (*RawLogParser, error) {
	inFile, err := os.Open(path)
	if err != nil {
		log.Panic(err.Error(), path)
		return nil, err
	}
	p := &RawLogParser{inputFile: inFile}
	return p, nil
}

func (p *RawLogParser) close() {
	p.inputFile.Close()
}
