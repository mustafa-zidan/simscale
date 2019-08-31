package main

import (
	"time"
)

type Logs []*Log

func (l Logs) Less(i, j int) bool {
	return l[i].Start.Before(l[j].Start)
}

func (l Logs) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l Logs) Len() int {
	return len(l)
}

type Log struct {
	Calls                        Logs
	Start, End                   time.Time
	Service, Span, Parent, Trace string
}

// Tree depth first insertion runs with
// O(log N) time complexity
func (l *Log) Insert(child *Log) bool {
	if l == nil {
		return false
	}
	if l.Span == child.Parent {
		Increment("success", 1)
		l.Calls = append(l.Calls, child)
		return true
	}
	for _, c := range l.Calls {
		if c.Insert(child) {

			return true
		}
	}

	return false
}

type LogTree struct {
	ID      string "json: id"
	Root    *Log   "json: root"
	Orphens map[string]Logs
}

//TODO check if maybe better to construct only on write
func (lt *LogTree) Insert(l *Log) {
	// create partial tree
	if calls, ok := lt.Orphens[l.Span]; ok {
		l.Calls = calls
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
		lt.Orphens[l.Parent] = Logs{l}
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
		Start:   start,
		End:     end,
		Service: properties["service"],
		Span:    properties["span"],
		Parent:  properties["caller_span"],
	}, nil
}
