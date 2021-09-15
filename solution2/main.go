package main

import (
	"fmt"
	"sync"
	"time"
)

const PHILOSOPHER_COUNT = 5
const ITERATIONS = 10

var cleanForks [PHILOSOPHER_COUNT]bool
var forkChannels [PHILOSOPHER_COUNT]chan bool
var forkOwners [PHILOSOPHER_COUNT]int
var wantsFork [PHILOSOPHER_COUNT]bool
var isEating [PHILOSOPHER_COUNT]bool

func main() {
	println("solution 2")
	startTime := time.Now()

	var waitGroup sync.WaitGroup

	for i := 0; i < PHILOSOPHER_COUNT; i++ {
		waitGroup.Add(1)
		forkChannels[i] = make(chan bool)

		// Setup fork owners in a way that prevents a
		// precedence cycle.
		forkOwners[i] = i
		if i == PHILOSOPHER_COUNT-1 {
			// Setting the last fork's owner to
			// the first philosopher breaks the cycle
			// that would form if we assigned a fork to
			// each philosopher.
			// There are many ways to break this cycle,
			// I just did it this way because it seemed the
			// simplest way to do it.
			forkOwners[i] = 0
		}
	}

	for philId := 0; philId < PHILOSOPHER_COUNT; philId++ {
		go runPhisolopher(philId, &waitGroup)
		go runForkMonitor(philId)
	}

	waitGroup.Wait()

	endTime := time.Now()

	fmt.Printf("Duration: %f\n", endTime.Sub(startTime).Seconds())
}

func runPhisolopher(philId int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	defer LogAction(philId, "Finished")

	for i := 0; i < ITERATIONS; i++ {
		LogAction(philId, "Thinking")
		think()

		leftFork := leftForkId(philId)
		rightFork := rightForkId(philId)

		// Wait for forks
		wantsFork[philId] = true

		LogAction(philId, "Waiting(L)")
		for forkOwners[leftFork] != philId {
			<-forkChannels[leftFork]
		}

		LogAction(philId, "Waiting(R)")
		for forkOwners[rightFork] != philId {
			<-forkChannels[rightFork]
		}

		LogAction(philId, "Eating")

		isEating[philId] = true
		wantsFork[philId] = false
		cleanForks[leftForkId(philId)] = false
		cleanForks[rightForkId(philId)] = false

		eat()

		isEating[philId] = false
	}
}

func runForkMonitor(forkId int) {
	for {
		tryGiveFork(forkId, leftPhilId(forkId))
		tryGiveFork(forkId, rightPhilId(forkId))
		time.Sleep(10 * time.Microsecond)
	}
}

func tryGiveFork(forkId int, philId int) {
	var otherPhilId int
	if leftPhilId(forkId) == philId {
		otherPhilId = rightPhilId(forkId)
	} else {
		otherPhilId = leftPhilId(forkId)
	}

	if wantsFork[philId] && (getForkPrecedence(forkId) == philId || !wantsFork[otherPhilId]) && !isEating[otherPhilId] {
		if forkOwners[forkId] != philId {
			cleanForks[forkId] = true
			forkOwners[forkId] = philId

			// notify a change has been made to the fork owner
			forkChannels[forkId] <- true
		}
	}
}

func eat() {
	time.Sleep(time.Millisecond * 100)
}

func think() {
	time.Sleep(time.Millisecond * 200)
}

func leftForkId(philId int) int {
	return (PHILOSOPHER_COUNT + philId - 1) % PHILOSOPHER_COUNT
}

func rightForkId(philId int) int {
	return philId
}

// Get the id of the philosopher that has precedence
// over a given fork
func getForkPrecedence(forkId int) int {
	leftPhil := leftPhilId(forkId)
	rightPhil := rightPhilId(forkId)
	forkOwner := forkOwners[forkId]
	isClean := cleanForks[forkId]

	// let u and v be left or right philosophers and u != v
	// u holds clean fork: u has precedence
	// u holds dirty fork: v has precedence
	if (leftPhil == forkOwner) == isClean {
		return leftPhil
	} else {
		return rightPhil
	}
}

func leftPhilId(forkId int) int {
	return forkId
}

func rightPhilId(forkId int) int {
	return (forkId + 1) % PHILOSOPHER_COUNT
}
