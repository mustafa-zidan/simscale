package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"
)

type Logs []*Log

func (l Logs) Less(i, j int) bool {
	return time.Time(l[i].Start).Before(time.Time(l[j].Start))
}

func (l Logs) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l Logs) Len() int {
	return len(l)
}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02T15:04:05.000Z07:00"))
	return []byte(stamp), nil
}

type Log struct {
	Start   JSONTime `json:"start"`
	End     JSONTime `json:"end"`
	Service string   `json:"service"`
	Span    string   `json:"span"`
	Calls   Logs     `json:"calls"`
	Parent  string   `json:"-"`
	Trace   string   `json:"-"`
}

// Tree depth first insertion runs with
// O(log N) time complexity
func (l *Log) Insert(child *Log) bool {
	if l == nil {
		return false
	}
	if l.Span == child.Parent {
		l.Calls = append(l.Calls, child)
		sort.Sort(l.Calls)
		return true
	}
	for _, c := range l.Calls {
		if c.Insert(child) {
			return true
		}
	}

	return false
}

func (l *Log) String() string {
	b, err := json.Marshal(l)
	if err != nil {
		log.Fatalf("Error: %s", err)
		return ""
	}
	return string(b)
}

type LogTree struct {
	ID      string          `json:"id"`
	Root    *Log            `json:"root"`
	Orphens map[string]Logs `json:"-"`
}

func (lt *LogTree) String() string {
	b, err := json.Marshal(lt)
	if err != nil {
		log.Fatalf("Error: %s", err)
		return ""
	}
	return string(b)
}

//TODO check if maybe better to construct only on write
func (lt *LogTree) Insert(l *Log) {
	// create partial tree
	if calls, ok := lt.Orphens[l.Span]; ok {
		l.Calls = calls
		sort.Sort(l.Calls)
		delete(lt.Orphens, l.Span)
	}
	if l.Parent == "null" {
		lt.Root = l
	} else if !lt.Root.Insert(l) {
		lt.AddOrphen(l)

	}
}

func (lt *LogTree) AddOrphen(l *Log) {

	if calls, ok := lt.Orphens[l.Parent]; ok {
		lt.Orphens[l.Parent] = append(calls, l)
	} else {
		//Look throught the current orphans for a parent
		var success = false
		for _, logs := range lt.Orphens {
			for _, ll := range logs {
				if ll.Insert(l) {
					success = true
					break
				}
			}
		}
		if !success {
			lt.Orphens[l.Parent] = Logs{l}
		}
	}
}

func NewLog(properties map[string]string) (*Log, error) {
	start, err := time.Parse(time.RFC3339, properties["start"])
	if err != nil {
		return nil, err
	}
	end, err := time.Parse(time.RFC3339, properties["end"])
	if err != nil {
		return nil, err
	}
	return &Log{
		Trace:   properties["trace"],
		Start:   JSONTime(start),
		End:     JSONTime(end),
		Service: properties["service"],
		Span:    properties["span"],
		Parent:  properties["caller_span"],
		Calls:   make(Logs, 0),
	}, nil
}
