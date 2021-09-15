package main

import (
	"fmt"
	"time"
)

const PHILOSOPHER_COUNT = 5
const ITERATIONS = 10

func main() {
	startTime := time.Now()

	for i := 0; i < ITERATIONS; i++ {
		for philId := 0; philId < PHILOSOPHER_COUNT; philId++ {
			think(philId)
		}
		for philId := 0; philId < PHILOSOPHER_COUNT; philId++ {
			eat(philId)
		}
	}

	fmt.Printf("Duration: %f\n", time.Now().Sub(startTime).Seconds())

}

func eat(id int) {
	LogAction(id, "Eating")
	time.Sleep(time.Millisecond * 100)
}

func think(id int) {
	LogAction(id, "Thinking")
	time.Sleep(time.Millisecond * 200)
}
