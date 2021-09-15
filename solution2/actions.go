package main

import (
	"fmt"
	"sync"
)

var logMutex sync.Mutex

func LogAction(philId int, action string) {
	logMutex.Lock()

	colWidth := 15
	tabCount := philId

	print("\033[1G")
	for i := 0; i < PHILOSOPHER_COUNT; i++ {
		fmt.Printf("\033[%dC|", colWidth-1)
	}

	print("\033[1G")
	fmt.Printf("\033[%dC %s\n", tabCount*colWidth, action)

	logMutex.Unlock()
}
