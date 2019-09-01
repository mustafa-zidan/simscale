package domain

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func rawLogs() []map[string]string {
	a := make([]map[string]string, 0)

	m1 := make(map[string]string)
	val := []byte(`{
		"start": "2013-10-23T10:12:37.129Z",
		"end": "2013-10-23T10:12:37.240Z",
		"trace": "nhdyl6hs",
		"service": "service6",
		"caller_span": "34z4ib4a",
		"span": "zgbtx32n" }`)
	json.Unmarshal(val, &m1)
	a = append(a, m1)

	m2 := make(map[string]string)
	val = []byte(`{
		"start": "2013-10-23T10:12:37.708Z",
		"end": "2013-10-23T10:12:37.724Z",
		"trace": "nhdyl6hs",
		"service": "service9",
		"caller_span": "34z4ib4a",
		"span": "ai67mto3" }`)
	json.Unmarshal(val, &m2)
	a = append(a, m2)

	m3 := make(map[string]string)
	val = []byte(`{
		"start": "2013-10-23T10:12:37.709Z",
		"end": "2013-10-23T10:12:37.715Z",
		"trace": "nhdyl6hs",
		"service": "service6",
		"caller_span": "ai67mto3",
		"span": "nxlpyoj7" }`)
	json.Unmarshal(val, &m3)
	a = append(a, m3)

	m4 := make(map[string]string)
	val = []byte(`{
		"start": "2013-10-23T10:12:37.127Z",
		"end": "2013-10-23T10:12:37.891Z",
		"trace": "nhdyl6hs",
		"service": "service7",
		"caller_span": "null",
		"span": "34z4ib4a" }`)
	json.Unmarshal(val, &m4)
	a = append(a, m4)
	return a
}

func TestLogInsertion(t *testing.T) {
	rawLogs := rawLogs()
	lt := &LogTree{ID: rawLogs[0]["trace"], Orphens: make(map[string]Logs)}
	for _, rl := range rawLogs {
		log.Println(rl)
		l, err := NewLog(rl)
		log.Println(l)
		assert.NotNil(t, l)
		assert.Nil(t, err)
		lt.Insert(l)
	}
	log.Println(lt)
	assert.NotNil(t, lt.Root)
	assert.Equal(t, 2, len(lt.Root.Calls))
	assert.Equal(t, 1, len(lt.Root.Calls[1].Calls))
	assert.Equal(t, "34z4ib4a", lt.Root.Span)
	assert.Equal(t, "zgbtx32n", lt.Root.Calls[0].Span)
	assert.Equal(t, "ai67mto3", lt.Root.Calls[1].Span)
	assert.Equal(t, "nxlpyoj7", lt.Root.Calls[1].Calls[0].Span)

}
