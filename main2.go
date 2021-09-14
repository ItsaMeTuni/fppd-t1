package main

import (
	"fmt"
	"sync"
	"time"
)

const PHILOSOPHER_COUNT = 5
const ITERATIONS = 10

var forkMutexes [PHILOSOPHER_COUNT]sync.Mutex
var grabbingMutex sync.Mutex

func main() {
	startTime := time.Now()

	var waitGroup sync.WaitGroup

	for philId := 0; philId < PHILOSOPHER_COUNT; philId++ {
		waitGroup.Add(1)
		go runPhisolopher(philId, &waitGroup)
	}

	waitGroup.Wait()

	endTime := time.Now()

	fmt.Printf("Duration: %f\n", endTime.Sub(startTime).Seconds())
}

func runPhisolopher(philId int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	for i := 0; i < ITERATIONS; i++ {
		LogAction(philId, "Thinking")
		think()

		leftFork := &forkMutexes[leftForkId(philId)]
		rightFork := &forkMutexes[rightForkId(philId)]

		// problem: basically creates a queue of philosofers wanting to grab
		// forks, even if the forks a philosopher wants are free. If a philosopher
		// wants a fork that's being used by someone else all other philosophers must
		// wait for that philosofer to get his fork in order to be able to get theirs,
		// even if they are different forks.
		LogAction(philId, "Waiting(G)")
		grabbingMutex.Lock()

		LogAction(philId, "Waiting(L)")
		leftFork.Lock()

		LogAction(philId, "Waiting(R)")
		rightFork.Lock()

		grabbingMutex.Unlock()

		LogAction(philId, "Eating")
		eat()

		leftFork.Unlock()
		rightFork.Unlock()
	}
}

func eat() {
	time.Sleep(time.Millisecond * 10)
}

func think() {
	time.Sleep(time.Millisecond * 20)
}

func leftForkId(philId int) int {
	return philId
}

func rightForkId(philId int) int {
	forkId := philId + 1

	if forkId >= PHILOSOPHER_COUNT {
		forkId -= PHILOSOPHER_COUNT
	}

	return forkId
}
