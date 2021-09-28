package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type AvailableBike struct {
	link string
	model string
}

func main() {
	fmt.Println("starting job at: ", time.Now())
	var bikesArray []AvailableBike
	htPage := getHtmlBody(os.Getenv("URL"))
	mvPage := getHtmlBody(os.Getenv("MAV_URL"))
	highTowerArray := findTheBikes(htPage, "hightower")
	maverickArray := findTheBikes(mvPage, "maverick")
	bikesArray = append(highTowerArray, maverickArray...)

	if(len(bikesArray) == 0) {
		fmt.Println("no entries were found")
		return
	}
	
	readDB(bikesArray)
	// seedDB(bikesAvailable)
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

func findTheBikes(bikeEntries *goquery.Document, modelName string) []AvailableBike {
	var bikesAvailable []AvailableBike
	bikeEntries.Find(".uImage").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Children().Attr("href")
		bikeFound := AvailableBike {
			model: modelName,
			link: link,
		}
		if(link != "https://www.pinkbike.com/buysell/3029717/") {
			bikesAvailable = append(bikesAvailable, bikeFound)
		}
	})
	return bikesAvailable
}
