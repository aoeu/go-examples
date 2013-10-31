// Download data from the MTA website and print out the train line statuses.
package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Line struct {
	Name   string `xml:"name"`
	Status string `xml:"status"`
}

type Service struct {
	Timestamp  string `xml:"timestamp"`
	Trainlines []Line `xml:"subway>line"`
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func downloadUrl(url string) []byte {
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	return body
}

func main() {
	data := downloadUrl("http://mta.info/status/serviceStatus.txt")
	service := new(Service)
	err := xml.Unmarshal(data, service)
	check(err)

	fmt.Println(service.Timestamp)
	for i, trainLine := range service.Trainlines {
		fmt.Printf("%d. Trainline %s is %s\n", i, trainLine.Name, trainLine.Status)
	}
}
