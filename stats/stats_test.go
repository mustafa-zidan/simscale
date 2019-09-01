package stats

import (
    "sync"
    "testing"
)

func benchmarkChannelsRoutine(b *testing.B, e chan bool) {
    for i := 0; i < b.N; i++ {
        Increment("abc123", 5)
        Increment("def456", 5)
        Increment("ghi789", 5)
        Increment("abc123", 5)
        Increment("def456", 5)
        Increment("ghi789", 5)
    }
    e <- true
}

func BenchmarkChannels(b *testing.B) {
    b.StopTimer()
    CounterInitialize()
    e := make(chan bool)
    b.StartTimer()

    go benchmarkChannelsRoutine(b, e)
    go benchmarkChannelsRoutine(b, e)
    go benchmarkChannelsRoutine(b, e)
    go benchmarkChannelsRoutine(b, e)
    go benchmarkChannelsRoutine(b, e)

    <-e
    <-e
    <-e
    <-e
    <-e

}

var mux sync.Mutex
var m map[string]int

func benchmarkMutexIncrement(bucket string, value int) {
    mux.Lock()
    m[bucket] += value
    mux.Unlock()
}

func benchmarkMutexRoutine(b *testing.B, e chan bool) {
    for i := 0; i < b.N; i++ {
        benchmarkMutexIncrement("abc123", 5)
        benchmarkMutexIncrement("def456", 5)
        benchmarkMutexIncrement("ghi789", 5)
        benchmarkMutexIncrement("abc123", 5)
        benchmarkMutexIncrement("def456", 5)
        benchmarkMutexIncrement("ghi789", 5)
    }
    e <- true
}

func BenchmarkMutex(b *testing.B) {
    b.StopTimer()
    m = make(map[string]int)
    e := make(chan bool)
    b.StartTimer()

    for i := 0; i < b.N; i++ {
        benchmarkMutexIncrement("abc123", 5)
        benchmarkMutexIncrement("def456", 5)
        benchmarkMutexIncrement("ghi789", 5)
        benchmarkMutexIncrement("abc123", 5)
        benchmarkMutexIncrement("def456", 5)
        benchmarkMutexIncrement("ghi789", 5)
    }

    go benchmarkMutexRoutine(b, e)
    go benchmarkMutexRoutine(b, e)
    go benchmarkMutexRoutine(b, e)
    go benchmarkMutexRoutine(b, e)
    go benchmarkMutexRoutine(b, e)

    <-e
    <-e
    <-e
    <-e
    <-e
}

// TODO Test counters working correcly
