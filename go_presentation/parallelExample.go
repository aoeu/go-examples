// This program reads URLs from a file, downloads them in parallel, and prints the titles.
package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func check(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func readFile(fileName string) []string {
	data, err := ioutil.ReadFile(fileName)
	check(err)
	lines := strings.Split(string(data), "\n")
	return lines[0 : len(lines)-1]
}

func scrapeTitle(url string, titles chan string, wg *sync.WaitGroup) {
	doc, err := goquery.NewDocument(url)
	check(err)
	titles <- doc.Find("#section_0").Text()
	wg.Done()
}

func downloadAsync(urls []string) {
	var wg sync.WaitGroup
	titles := make(chan string, len(urls))
	for _, url := range urls {
		wg.Add(1)
		go scrapeTitle(url, titles, &wg)
	}
	wg.Wait()
	close(titles)
	for title := range titles {
		fmt.Println(title)
	}
}

func downloadSync(urls []string) {
	for _, url := range urls {
		doc, err := goquery.NewDocument(url)
		check(err)
		fmt.Println(doc.Find("#section_0").Text())
	}
}

func main() {
	urls := readFile("urls.txt")

	start := time.Now()
	downloadSync(urls)
	end := time.Now()
	fmt.Println(end.Sub(start))

	start = time.Now()
	downloadAsync(urls)
	end = time.Now()
	fmt.Println(end.Sub(start))
}
