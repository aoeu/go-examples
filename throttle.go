package main

import (
	"fmt"
	"time"
)

func main() {
	maxGoRoutines := 2
	throttle := make(chan bool, maxGoRoutines)
	wordsToPrint := []string{"foo", "bar", "baz", "qux", "zip", "ding", "poof"}

	for _, word := range wordsToPrint {
		throttle <- true // Blocks until we can send to the channel.
		go func(word string) {
			defer func() { <-throttle }() // Signal that this go routine has finished.
			time.Sleep(1 * time.Second)
			fmt.Println(word)
		}(word)
	}

	// The below loop will block until all go routines are finished.
	for i := 0; i < cap(throttle); i++ {
		throttle <- true
	}
}
