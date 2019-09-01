package cache

import (
	"github.com/ReneKroon/ttlcache"
	. "github.com/mustafa-zidan/simscale/domain"
	"github.com/mustafa-zidan/simscale/io"
	"github.com/mustafa-zidan/simscale/stats"
	"sync"
	"time"
)

type Cache struct {
	writer       *io.FileWriter
	readChannel  chan *readLogTreeOp
	writeChannel chan *Log
	stopChannel  chan bool
	wg           *sync.WaitGroup
	mux          sync.Mutex
}

// struct which is used communicate over the read channel
type readLogTreeOp struct {
	Trace   string
	LogTree chan *LogTree
}

func (c *Cache) writeToFile(key string, value interface{}) {
	c.mux.Lock()
	defer c.wg.Done()
	lt := value.(*LogTree)
	if lt.Root != nil {
		c.writer.Append(lt)
	}
	stats.Increment("orphens", len(lt.Orphens))
	c.mux.Unlock()
}

func NewCache(filePath string, wg *sync.WaitGroup) *Cache {
	c := &Cache{
		writer:       io.NewFileWriter(filePath),
		readChannel:  make(chan *readLogTreeOp),
		writeChannel: make(chan *Log),
		stopChannel:  make(chan bool),
		wg:           wg,
	}
	go c.goCacheWriter()

	return c
}

func (c *Cache) goCacheWriter() {
	cache := ttlcache.NewCache()
	cache.SetTTL(time.Duration(1 * time.Second))
	cache.SetExpirationCallback(c.writeToFile)

	for {
		select {
		// read from the current state map
		case read := <-c.readChannel:
			if value, exists := cache.Get(read.Trace); exists {
				read.LogTree <- value.(*LogTree)
			}
		case insert := <-c.writeChannel:
			c.mux.Lock()
			if value, exists := cache.Get(insert.Trace); exists {
				tree := value.(*LogTree)
				tree.Insert(insert)
			} else {
				c.wg.Add(1)
				tree := &LogTree{ID: insert.Trace, Orphens: make(map[string]Logs)}
				tree.Insert(insert)
				cache.Set(insert.Trace, tree)
			}
			stats.Increment("success", 1)
			c.mux.Unlock()
		case <-c.stopChannel: // STOP
			close(c.readChannel)
			close(c.writeChannel)
			close(c.stopChannel)
			break
		}
	}
	c.writer.Close()
}

func (c *Cache) Insert(l *Log) {
	c.writeChannel <- l
}
