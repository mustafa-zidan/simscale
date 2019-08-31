package main

type Counter struct {
	Bucket string
	Value  int
}

type CounterQuery struct {
	Bucket  string
	Channel chan int
}

var counter map[string]int
var counterChan chan Counter
var counterQueryChan chan CounterQuery
var counterListChan chan chan map[string]int

func CounterInitialize() {
	counter = make(map[string]int)
	counterChan = make(chan Counter, 0)
	counterQueryChan = make(chan CounterQuery, 100)
	counterListChan = make(chan chan map[string]int, 100)
	go goCounterWriter()
}

func goCounterWriter() {
	for {
		select {
		case ci := <-counterChan:
			if len(ci.Bucket) == 0 {
				return
			}
			counter[ci.Bucket] += ci.Value
			break
		case cq := <-counterQueryChan:
			val, found := counter[cq.Bucket]
			if found {
				cq.Channel <- val
			} else {
				cq.Channel <- -1
			}
			break
		case cl := <-counterListChan:
			nm := make(map[string]int)
			for k, v := range counter {
				nm[k] = v
			}
			cl <- nm
			break
		}
	}
}

func Increment(bucket string, counter int) {
	if len(bucket) == 0 || counter == 0 {
		return
	}
	counterChan <- Counter{bucket, counter}
}

func GetCounter(bucket string) int {
	if len(bucket) == 0 {
		return -1
	}
	reply := make(chan int)
	counterQueryChan <- CounterQuery{bucket, reply}
	return <-reply
}

func CounterList() map[string]int {
	reply := make(chan map[string]int)
	counterListChan <- reply
	return <-reply
}
