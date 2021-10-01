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

type BikesToLookFor struct {
	urlBodyPage *goquery.Document
	modelPage string
}
var bikesAvailable []AvailableBike

func main() {
	fmt.Println("starting job at: ", time.Now())
	var bikesArray []AvailableBike
	
	bikesToSearch := []BikesToLookFor {
		{ urlBodyPage: getHtmlBody(os.Getenv("URL")), modelPage: "htPage", }, 
		{ urlBodyPage: getHtmlBody(os.Getenv("MAV_URL")), modelPage: "mvPage", },
		{ urlBodyPage: getHtmlBody(os.Getenv("TRUL")), modelPage: "stPage", },
	}
	
	for _, bikePage := range bikesToSearch {
		bikePageArray := findTheBikes(bikePage.urlBodyPage, bikePage.modelPage)
		bikesArray = append(bikePageArray)
	}

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
