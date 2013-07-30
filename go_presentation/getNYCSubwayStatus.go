package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Line struct {
	Name   string `xml:"name"`
	Status string `xml:"status"`
}

type ServiceStatus struct {
	Timestamp  string `xml:"timestamp"`
	TrainLines []Line `xml:"subway>line"`
}

func check(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func downloadUrl(url string) []byte {
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	check(err)
	return data
}

func main() {
	data := downloadUrl("http://mta.info/status/serviceStatus.txt")
	status := new(ServiceStatus)
	err := xml.Unmarshal(data, status)
	check(err)
	fmt.Println(status.Timestamp)
	for _, trainLine := range status.TrainLines {
		fmt.Printf("Train line: %s - %s\n", trainLine.Name, trainLine.Status)
	}
}
