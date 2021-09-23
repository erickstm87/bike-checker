package main

import (
	"fmt"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"os"
)

func main() {
	bikePage := getHtmlBody("https://www.pinkbike.com/buysell/list/?region=3&q=hightower&framesize=9,11,12,17,18,20,21,22")
	found := false
	
	bikePage.Find(".uImage").Each(func(i int, s *goquery.Selection) {
		found = true
		link, _ := s.Children().Attr("href")
		fmt.Println(link)
	})
	if(!found) {
		fmt.Println("no entries were found")
		os.Exit(1)
	}
}

func errorHandler(err error) {
	fmt.Println("there was an error: ", err)
}

func getHtmlBody(url string) *goquery.Document {
	resp, err := http.Get(url)
	if err != nil {
		errorHandler(err)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		errorHandler(err)
	} 
	return doc
}

// https://www.pinkbike.com/buysell/list/?region=3&q=hightower&framesize=9,11,12,17,18,20,21,22,16,19,24,25,26,28,29
// https://www.pinkbike.com/buysell/list/?region=3&q=hightower&framesize=9,11,12,17,18,20,21,22