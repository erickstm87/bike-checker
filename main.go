package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	var parsedArray []string
	fmt.Println("starting job at: ", time.Now())
	bikePage := getHtmlBody(os.Getenv("URL"))
	found := false
	bikesAvailable := []string{""}
	
	bikePage.Find(".uImage").Each(func(i int, s *goquery.Selection) {
		found = true
		link, _ := s.Children().Attr("href")
		bikesAvailable = append(bikesAvailable, link)
	})
	if(!found) {
		fmt.Println("no entries were found")
		return
	}
	// erase blank spaces from array by copying to new
	for i := range bikesAvailable {
		if(bikesAvailable[i] != "" && bikesAvailable[i] != " ") {
			parsedArray = append(parsedArray, bikesAvailable[i])
		}
	}
	readDB(bikesAvailable)
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
