# Overview
Overview of how the app operates

## Datastructure
The datastructure for this coding task is prettry straight forward which resembles
N-ary tree.

```go
type Logs []*Log

type LogTree struct {
    ID      string          `json:"id"`
    Root    *Log            `json:"root"`
    Orphens map[string]Logs `json:"-"`
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
```
`LogTree` is the parent of all the logs that contains the same trace with a pointer to the `Root` `Log`
and since the logs comes out of order we need to have temporary storage for logs that are still has no
parent added to the tree yet.

`Orphens` map holds List of partial trees so on the insertion call in the `LogTree` we try to do the following:
- Check if the current log has orphens then we attach them to the current log to extend the partial tree
- If the current Log is Root which is determaned if the Parent of the log is `null` then place it in the `LogTree` Root
- Else if the `Root` already exists then we perform a simple tree insertion with `O(log N)` time complexity
- If it could not be inserted due to missing parent then we add it to the map of Orphans and this is done by
    - Try to insert the log in each of the orphans since we might have the parent as part of the aub tree of another Orphen
    - If it does not exist then create a entry for the log parent

### JSON Conversion
JSON conversion is done using `encoding/json` built-in go library by setting `json:"<name_mapping>"` tag on the structure field.

> **NOTE:** I need to create a new type called JSONTime due to discrepancy in nanoseconds format between output
Examples and the output here for example on the output examples attached time is `2013-10-23T10:12:37.240Z` but the
output here is `2013-10-23T10:12:37.24Z`

```go
type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
    stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02T15:04:05.000Z07:00"))
    return []byte(stamp), nil
}
```

## IO
I needed to create wrapper for `FileReader` and `FileWriter` under `io` directory to try to decouple all IO operations
from the `Parser` and the `Cache`to make it more testable.

I Created Interfaces for Both `Reader` and `Writer` that look like the following:

```go
type Reader interface {
    Open()
    Read(chan string)
    Close()
}

type Writer interface {
    Append(*LogTree)
    Open(string)
    Close()
}
```

 The reason for that is to be able to create `Mocks` that implements the same methods and pass them to Both `Parser` and
 `Cache`

 ### File Reader
Reading File leveraging goroutines and channel to be able to process raw log lines concurrently

In the Parser File Reading is done as following:

```go
defer p.reader.Close()

lines := make(chan string)
go p.reader.Read(lines)

for line := range lines {
    //Parse and Insert Lines to LogTrees in the Cache
}
```

### File Writer
File Writer using `bufio.Writer` to reduce IO operations as much as possible
```go
func (f *FileWriter) Append(logtree *LogTree) {
    lt, err := json.Marshal(logtree)
    if err != nil {
        log.Fatalf("failed marshalling logtree to JSON: %s", err)
    }
    f.writer.WriteString(fmt.Sprintf("%s\n", lt))
}
```
> **NOTE:** Due to using channels and goroutines in the cache I have to use Mutex to lock the writer otherwise the output gets all jumbled up

## Cache
In Caching since the logs come out of order and I didn't want to load the whole file in memory while the processing is complete
I used an external Lightweight Library [`ttlcache`](https://github.com/ReneKroon/ttlcache) which Mixture between LRU and TTL cache
- It uses `PriortryQueue` datastructure to keep alive all active Items
- Perform eviction only when TTL passed without any operation or activity on the item.

To acheive High Concurrency on The cache I setup verious channels for reading and writing items to/from the cache.

The cache Uses `sync.WaitGroup` to ensure the app waits till all items are evicted and written into file.
I increment the waiting group by one on each LogTree Creation and decrement it only on an eviction of LogTree from the cach

On Writing to the Cache
``` go
...
if value, exists := cache.Get(insert.Trace); exists {
    tree := value.(*LogTree)
    tree.Insert(insert)
} else {
    c.wg.Add(1) //Only increment waiting group in the case of creating new LogTree
    tree := &LogTree{ID: insert.Trace, Orphens: make(map[string]Logs)}
    tree.Insert(log)
    cache.Set(insert.Trace, tree)
}
...
```
On Eviction
```go
...
defer c.wg.Done() //On Cache Item Eviction
lt := value.(*LogTree)
if lt.Root != nil {
    c.writer.Append(lt)
}
...
```

### Limitations in LRU/TTL Cache
Since [`ttlcache`](https://github.com/ReneKroon/ttlcache) is a generic Library; Getting all Items from the Cache is not Possible
which means I need to wait all the items to be evicted on their own and not to flush the whole cache at once.

And even if this functionality is provided it would return a `map[string]interface{}` so I needed to loop over the items and
do inline type assertion
```go
lt := item.(*LogTree)
```
which render the whole idea useless.

## Parser
Parser is responsible of validating the format of the raw log and create the `Log` entry to be passed to the cache.

Parser used `regexp` to create named submatches using the following expression
```go
pattern  := "(?P<start>.*)\\s(?P<end>.*)\\s(?P<trace>.*)\\s(?P<service>.*)\\s(?P<caller_span>.*)->(?P<span>.*)"
re       := regexp.MustCompile(pattern)
expNames := re.SubexpNames()
```
Then extract the named submatches to a map and create a log using the mapped properties.

```go
func (p *RawLogParser) Parse(line string) (*Log, error) {
    matches := re.FindAllStringSubmatch(line, -1)
    if len(matches) == 0 {
        return nil, fmt.Errorf("Malformed Log Format %s, Skipping!", line)
    }

    m := map[string]string{}
    for i, n := range matches[0] {
        m[expNames[i]] = n
    }
    return NewLog(m)
}
```


## Stats
Using channels, go routines and the same mechanism used in caching but simpler I was able to collect the following information
- Total Number of Logs
- Number of Records parsed
- Number of Traces Created
- Number of Orphen Logs
- Number of Malformed Logs

The underlaying structure of stats is a simple Key/Value datastructure

``` Go
type Counter struct {
    Bucket string
    Value  int
}
```

Then using `stats.Increment(key, value)` function to we add values to a certain key

``` go
// Parse takes a line then parse it using regex to Log Object
func (p *RawLogParser) Parse(line string) (*Log, error) {
    ...
    if len(matches) == 0 {
        stats.Increment("malformed", 1)
        return nil, fmt.Errorf("Malformed Log Format %s, Skipping!", line)
    }
    ...
}
```

Then using `stats.CounterList()` to get the map of stats accumulated.

```go
func DoSomething() {
    ...
    s := stats.CounterList()
    ...
}

```

## Final Notes:
- Calculating more stats is possible - average tree depth etc.- but would be more complex and maybe require significant change.
- In production mode I would add like to have more tests in place.
