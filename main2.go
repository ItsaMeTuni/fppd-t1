package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type Action struct {
	name      string
	timestamp time.Time
}

const PHILOSOPHER_COUNT = 5
const ITERATIONS = 10

var forkMutexes [PHILOSOPHER_COUNT]sync.Mutex
var grabbingMutex sync.Mutex

var actions [][]Action

func main() {
	startTime := time.Now()
	actions = make([][]Action, PHILOSOPHER_COUNT)

	var waitGroup sync.WaitGroup

	for philId := 0; philId < PHILOSOPHER_COUNT; philId++ {
		actions[philId] = make([]Action, 0)

		waitGroup.Add(1)
		go runPhisolopher(philId, &waitGroup)
	}

	waitGroup.Wait()

	endTime := time.Now()

	fmt.Printf("Duration: %f\n", endTime.Sub(startTime).Seconds())

	printTimings(endTime.Sub(startTime))
}

func runPhisolopher(philId int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	for i := 0; i < ITERATIONS; i++ {
		actions[philId] = append(actions[philId], Action{name: "Thinking", timestamp: time.Now()})
		think()

		leftFork := &forkMutexes[leftForkId(philId)]
		rightFork := &forkMutexes[rightForkId(philId)]

		// problem: basically creates a queue of philosofers wanting to grab
		// forks, even if the forks a philosopher wants are free. If a philosopher
		// wants a fork that's being used by someone else all other philosophers must
		// wait for that philosofer to get his fork in order to be able to get theirs,
		// even if they are different forks.
		actions[philId] = append(actions[philId], Action{name: "Waiting(G)", timestamp: time.Now()})
		grabbingMutex.Lock()
		actions[philId] = append(actions[philId], Action{name: "Waiting(L)", timestamp: time.Now()})
		leftFork.Lock()
		actions[philId] = append(actions[philId], Action{name: "Waiting(R)", timestamp: time.Now()})
		rightFork.Lock()
		grabbingMutex.Unlock()

		actions[philId] = append(actions[philId], Action{name: "Eating", timestamp: time.Now()})
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

type IdentifiedAction struct {
	Action
	philId int
}

type ByTimestamp []IdentifiedAction

func (a ByTimestamp) Len() int      { return len(a) }
func (a ByTimestamp) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool {
	return a[i].timestamp.Nanosecond() < a[j].timestamp.Nanosecond()
}

func printTimings(duration time.Duration) {
	var allActions []IdentifiedAction = make([]IdentifiedAction, 0)
	for philId := 0; philId < PHILOSOPHER_COUNT; philId++ {
		for _, action := range actions[philId] {
			allActions = append(allActions, IdentifiedAction{
				philId: philId,
				Action: action,
			})
		}
	}

	sort.Sort(ByTimestamp(allActions))

	colWidth := 15

	for i := 0; i < PHILOSOPHER_COUNT; i++ {
		fmt.Printf("\033[%dC|", colWidth-1)
	}

	for i := 0; i < PHILOSOPHER_COUNT; i++ {
		fmt.Printf("\033[1G\033[%dC%d", i*colWidth, i+1)
	}

	print("\n")

	for _, action := range allActions {
		tabCount := action.philId

		print("\033[1G")
		for i := 0; i < PHILOSOPHER_COUNT; i++ {
			fmt.Printf("\033[%dC|", colWidth-1)
		}

		print("\033[1G")
		fmt.Printf("\033[%dC %s\n", tabCount*colWidth, action.name)
	}

}
