package main

import (
	"github.com/ReneKroon/ttlcache"
	"log"
	"time"
)

// struct which is used communicate over the read channel
type readLogTreeOp struct {
	Trace   string
	LogTree chan *LogTree
}

var readLogTreeChannel chan *readLogTreeOp
var insertLogChannel chan *Log
var stopCacheChannel chan bool

func newItemCallback(key string, value LogTree) {
	log.Printf("New key(%s) added\n", key)
}

func checkExpirationCallback(key string, value LogTree) bool {
	if key == "key1" {
		// if the key equals "key1", the value
		// will not be allowed to expire
		return false
	}
	// all other values are allowed to expire
	return true
}

func expirationCallback(key string, value interface{}) {
	//TODO write to file and add orphans count to stats
	log.Printf("This key(%s) has expired\n", key)
}

func InitTTLCache() {
	readLogTreeChannel = make(chan *readLogTreeOp)
	insertLogChannel = make(chan *Log)
	stopCacheChannel = make(chan bool)

	go goCacheWriter()
}

func goCacheWriter() {
	cache := ttlcache.NewCache()
	cache.SetTTL(time.Duration(10 * time.Second))
	cache.SetExpirationCallback(expirationCallback)

	for {
		select {
		// read from the current state map
		case read := <-readLogTreeChannel:
			if value, exists := cache.Get(read.Trace); exists {
				read.LogTree <- value.(*LogTree)
			}
		case insert := <-insertLogChannel:

			if value, exists := cache.Get(insert.Trace); exists {
				tree := value.(*LogTree)
				tree.Insert(insert)
			} else {
				tree := &LogTree{ID: insert.Trace, Orphens: make(map[string]Logs)}
				tree.Insert(insert)
				cache.Set(insert.Trace, tree)
			}

		case <-stopCacheChannel: // STOP

			close(readLogTreeChannel)
			close(insertLogChannel)
			break
		}
	}
}

func InsertLogToCache(l *Log) {
	insertLogChannel <- l
}
