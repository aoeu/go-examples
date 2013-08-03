package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/xml"
	"github.com/davecgh/go-spew/spew"
)

type RSS struct {
	Version string `xml:"version,attr"`
	// Using shorthand of just "Channel" doesn't work here, e.g.
	// Channel `xml:"channel"`
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Link string `xml:"link"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title string `xml:"title"`
	Link string `xml:"link"`
	Description string `xml:"description"`
	Guid string `xml:"guid"`
	PubDate string `xml:"pubDate"`
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func download(url string) []byte {
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	check(err)
	return data
}

func main() {
	data := download("http://xkcd.com/rss.xml")
	rss := new(RSS)
	xml.Unmarshal(data, rss)
	log.Println(len(rss.Channel.Items))
	spew.Dump(rss)
}
