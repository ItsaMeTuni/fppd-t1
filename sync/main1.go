package main

import (
	"fmt"
	"time"
)

func main() {
	startTime := time.Now()

	for philId := 0; philId < 5; philId++ {
		runPhisolopher(philId)
	}

	fmt.Printf("Duration: %f\n", time.Now().Sub(startTime).Seconds())

}

func runPhisolopher(id int) {
	for i := 0; i < 10; i++ {
		fmt.Printf("%d thinking...\n", id)
		think()
		fmt.Printf("%d eating...\n", id)
		eat()
	}
}

func eat() {
	time.Sleep(time.Millisecond * 100)
}

func think() {
	time.Sleep(time.Millisecond * 200)
}
