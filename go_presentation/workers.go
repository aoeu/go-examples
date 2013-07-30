// This program demonstrates how to throttle tasks to a finite number of workers
// while reporting results to a manager.
package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"time"
)

type Worker struct {
	Id               int
	Results          chan Result
	AvailableWorkers chan Worker
}

func (w Worker) Scrape(url string) {
	hostname := "https://en.m.wikipedia.org/"
	doc, _ := goquery.NewDocument(url)
	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		headerText := tr.Find("th").Text()
		if headerText == "Influenced" {
			tr.Find("a").Each(func(i int, a *goquery.Selection) {
				pageName, _ := a.Attr("href")
				w.Results <- Result{time.Now().Unix(), hostname + pageName}
			})
		}
	})
}

type Result struct {
	TimeFinished int64 // Unix timestamp
	Payload      string
}

type Manager struct {
	Quit             chan bool
	Results          chan Result
	UrlsToDownload   chan string
	AvailableWorkers chan Worker
	TimeoutLength    int64
	LatestResultTime int64
	Downloaded       map[string]bool
}

func NewManager(numWorkers int, timeoutLength int64) *Manager {
	m := new(Manager)
	m.Quit = make(chan bool, 1)
	m.Results = make(chan Result, numWorkers)
	m.UrlsToDownload = make(chan string)
	m.TimeoutLength = timeoutLength
	m.AvailableWorkers = make(chan Worker, numWorkers)
	m.Downloaded = make(map[string]bool)
	for i := 0; i < numWorkers; i++ {
		m.AvailableWorkers <- Worker{i, m.Results, m.AvailableWorkers}
	}
	return m
}

func (m *Manager) Manage() {
	go m.Timeout()
	for {
		select {
		case <-m.Quit:
			m.Quit <- true
			return
		case r := <-m.Results:
			go m.Process(r)
		case url := <-m.UrlsToDownload:
			go m.Dispatch(url)

		}
	}
}

func (m *Manager) Process(r Result) {
	m.LatestResultTime = r.TimeFinished
	log.Println(r.Payload)
	if _, ok := m.Downloaded[r.Payload]; !ok {
		m.UrlsToDownload <- r.Payload
	}
}

func (m *Manager) Dispatch(url string) {
	select {
	case w := <-m.AvailableWorkers:
		go w.Scrape(url)
		m.Downloaded[url] = true
	case <-m.Quit:
		m.Quit <- true
		return
	}
}

func (m *Manager) Timeout() {
	m.LatestResultTime = time.Now().Unix()
	for {
		if time.Now().Unix()-m.LatestResultTime > m.TimeoutLength {
			log.Println("Timing out.")
			m.Quit <- true
			return
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func main() {
	m := NewManager(3, 5)
	go m.Manage()
	m.UrlsToDownload <- "https://en.m.wikipedia.org/wiki/C_programming_language"
	<-m.Quit
	log.Println("Done.")
}
