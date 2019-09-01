package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testNewParser(t *testing.T) {
	p, err := NewParser("test.json")
	assert.NotNil(t, err)
	assert.Nil(t, p)
}

func BenchmarkRegexParsingAsync(b *testing.B) {
	b.StopTimer()
	p := &RawLogParser{}
	e := make(chan bool)
	b.StartTimer()

	go benchmarkRegexParsingRoutine(b, e, p)
	go benchmarkRegexParsingRoutine(b, e, p)
	go benchmarkRegexParsingRoutine(b, e, p)
	go benchmarkRegexParsingRoutine(b, e, p)
	go benchmarkRegexParsingRoutine(b, e, p)

	<-e
	<-e
	<-e
	<-e
	<-e
}

func benchmarkRegexParsingRoutine(b *testing.B, e chan bool, p *RawLogParser) {
	for i := 0; i < b.N; i++ {
		p.Parse("2013-10-23T10:12:35.298Z 2013-10-23T10:12:35.300Z eckakaau service3 d6m3shqy->62d45qeh")
	}
	e <- true
}

func BenchmarkRegexParsing(b *testing.B) {
	p := &RawLogParser{}
	for i := 0; i < b.N; i++ {
		p.Parse("2013-10-23T10:12:35.298Z 2013-10-23T10:12:35.300Z eckakaau service3 d6m3shqy->62d45qeh")
	}
}
