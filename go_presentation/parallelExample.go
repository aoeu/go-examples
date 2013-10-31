package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	urls := readFile("urls.txt")

	fmt.Println("Synchronous downloading...")
	start := time.Now()
	downloadUrlsSync(urls)
	fmt.Println(time.Now().Sub(start))

	fmt.Println("Async...")
	start = time.Now()
	downloadUrls(urls)
	fmt.Println(time.Now().Sub(start))
}

func readFile(fileName string) (lines []string) {
	data, err := ioutil.ReadFile(fileName)
	check(err)
	lines = strings.Split(string(data), "\n")
	return lines[0 : len(lines)-1]
}

func downloadUrlsSync(urls []string) {
	for _, url := range urls {
		doc, err := goquery.NewDocument(url)
		check(err)
		fmt.Println(doc.Find("#section_0").Text())
	}
}

func downloadUrls(urls []string) {
	var wg sync.WaitGroup
	titles := make(chan string, len(urls))
	for _, url := range urls {
		wg.Add(1)
		go scrape(url, titles, &wg)
	}
	wg.Wait()
	close(titles)
	for title := range titles {
		fmt.Println(title)
	}
}

func scrape(url string, titles chan string, wg *sync.WaitGroup) {
	doc, err := goquery.NewDocument(url)
	check(err)
	titles <- doc.Find("#section_0").Text()
	wg.Done()
}
