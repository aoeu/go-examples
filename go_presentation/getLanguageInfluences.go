package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
)

func check(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func getInfluences(hostname, root string) {
	doc, err := goquery.NewDocument(hostname + root)
	check(err)
	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		headerText := tr.Find("th").Text()
		if headerText == "Influenced" {
			tr.Find("a").Each(func(i int, a *goquery.Selection) {
				pageName, _ := a.Attr("href")
				fmt.Println(hostname + pageName)
			})
		}
	})
}

func main() {
	hostname := "http://en.m.wikipedia.org/wiki"
	getInfluences(hostname, "/C_programming_language")
}
