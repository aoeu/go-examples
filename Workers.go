package main

import (
	"fmt"
	"time"
)

type Manager struct {
	FreeWorkers chan *Worker
	Tasks       chan int
	Results     chan int64
	LastResultTime int64 // Unix timestamp.
	TimeoutLength int64 // Number of seconds.
	Quit        chan bool
}

func NewManager(numWorkers int, timeoutLength int64) Manager {
	m := Manager{}
	m.TimeoutLength = timeoutLength
	m.LastResultTime = time.Now().Unix()
	m.FreeWorkers = make(chan *Worker, numWorkers)
	m.Tasks = make(chan int)
	m.Results = make(chan int64)
	m.Quit = make(chan bool, 1)
	return m
}

func (m *Manager) Manage() {
	for {
		select {
		case task := <-m.Tasks:
			go m.Dispatch(task)
		case result := <-m.Results:
			go m.Process(result)
		case <-m.Quit:
			m.Quit <- true
			return
		}
	}
}

func (m *Manager) Dispatch(task int) {
	fmt.Println("Dispatching task:", task)
	select {
	case worker := <-m.FreeWorkers:
		go worker.Work(task)
	case <-m.Quit:
		m.Quit <- true
		return
	}
}

func (m *Manager) Process(result int64) {
	fmt.Println("Got result at:", result)
	m.LastResultTime = int64(result)
}

func (m *Manager) Timeout(seconds int) {
	for {
		if time.Now().Unix() - m.LastResultTime > int64(seconds) {
			fmt.Println("Timing out.")
			m.Quit <- true
			return
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

type Worker struct {
	Id          int
	Results     chan int64
	FreeWorkers chan *Worker
}

func (w *Worker) Work(task int) {
	fmt.Printf("Worker %d is working on task %d\n", w.Id, task)
	time.Sleep(2 * time.Second)
	w.Results <- time.Now().Unix()
	w.FreeWorkers <- w    // put self on the shared availability "queue"
}

func main() {
	// Initialize.
	numWorkers := 3
	manager := NewManager(numWorkers, 3)
	for i := 0; i < numWorkers; i++ { /
		manager.FreeWorkers <- &Worker{i, manager.Results, manager.FreeWorkers}
	}
	// Run concurrently.
	go manager.Manage()  
	for i := 1; i <= 10; i++ {
		manager.Tasks <- i
	}
	// Wait.
	go manager.Timeout(3)
	<-manager.Quit

	fmt.Println("Done.")
}
