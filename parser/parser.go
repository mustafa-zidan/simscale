package parser

import (
	"fmt"
	"github.com/mustafa-zidan/simscale/cache"
	. "github.com/mustafa-zidan/simscale/domain"
	"github.com/mustafa-zidan/simscale/io"
	"github.com/mustafa-zidan/simscale/stats"
	"log"
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
	reader *io.FileReader
	cache  *cache.Cache
}

func (p *RawLogParser) Process() error {
	defer p.reader.Close()
	lines := make(chan string)

	go p.reader.Read(lines)
	for line := range lines {
		//TODO: Add counter lines processed
		stats.Increment("total", 1)
		l, err := p.Parse(line) // the line
		if err != nil {
			log.Println(err)
		}

		p.cache.Insert(l)
	}
	return nil
}

// Parse takes a line then parse it using regex to Log Object
func (p *RawLogParser) Parse(line string) (*Log, error) {
	matches := re.FindAllStringSubmatch(line, -1)
	if len(matches) == 0 {
		stats.Increment("ignored", 1)
		return nil, fmt.Errorf("Invalid Log Format %s, Skipping!", line)
	}

	m := map[string]string{}
	for i, n := range matches[0] {
		m[expNames[i]] = n
	}
	return NewLog(m)
}

func NewParser(filePath string, cache *cache.Cache) (*RawLogParser, error) {
	fileReader := io.NewFileReader(filePath)
	p := &RawLogParser{reader: fileReader, cache: cache}
	return p, nil
}
