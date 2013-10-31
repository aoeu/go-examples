package main

import (
	"fmt"
	"time"
)

func printWord(word string, throttle chan bool) {
	rateLimit := time.After(1 * time.Second)
	defer func() {
		<-rateLimit // Blocks until 1 second has passed.
		<-throttle  // Signal we've finished.
	}()
	fmt.Println(word)
}

func main() {
	maxGoRoutines := 2
	throttle := make(chan bool, maxGoRoutines)
	wordsToPrint := []string{"foo", "bar", "baz", "qux", "zip", "ding", "poof"}

	for _, word := range wordsToPrint {
		throttle <- true // Blocks until we can send to the channel.
		go printWord(word, throttle)
	}
	// Block until all go routines are finished.
	for i := 0; i < cap(throttle); i++ {
		throttle <- true
	}
}
