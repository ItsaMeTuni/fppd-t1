package main

import (
	"sync"
	"time"
)

type TimedMutex struct {
	sync.Mutex
	waitingTimestamps  map[int][]int
	acquiredTimestamps map[int][]int
	releasedTimestamps map[int][]int
}

func (tm *TimedMutex) Lock(philId int) {
	if tm.waitingTimestamps == nil {
		tm.waitingTimestamps = make(map[int][]int)
		tm.acquiredTimestamps = make(map[int][]int)
		tm.releasedTimestamps = make(map[int][]int)
	}

	// waiting
	tm.waitingTimestamps[philId] = append(tm.waitingTimestamps[philId], time.Now().Nanosecond())

	tm.Mutex.Lock()

	// acquired
	tm.acquiredTimestamps[philId] = append(tm.acquiredTimestamps[philId], time.Now().Nanosecond())
}

func (tm *TimedMutex) Unlock(philId int) {
	tm.releasedTimestamps[philId] = append(tm.releasedTimestamps[philId], time.Now().Nanosecond())

	tm.Mutex.Unlock()
}
