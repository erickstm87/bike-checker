package main

import (
	"fmt"
	"net/http"
	"time"
	"strconv"
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
	loc, _ := time.LoadLocation("MST")
	theTime := time.Now().In(loc)
	year := theTime.Year()
	month := theTime.Month()
	day := theTime.Day()
	hour := theTime.Hour()
	minute := theTime.Minute()
	sMinute := ""
	if(theTime.Minute() < 10) {
		sMinute = "0" + strconv.Itoa(minute)
	} else {
		sMinute = strconv.Itoa(minute)
	}
    
	theDate := strconv.Itoa(year) + "-" + month.String() + "-" + strconv.Itoa(day) + " " + strconv.Itoa(hour) + ":" + sMinute

	fmt.Println("starting job at: ", theDate)
	var bikesArray []AvailableBike
	bikesToSearch := []BikesToLookFor {
		{ urlBodyPage: getHtmlBody("https://www.pinkbike.com/buysell/list/?region=3&q=hightower&framesize=9,11,12,17,18,20,21,22"), modelPage: "hightower", }, 
		{ urlBodyPage: getHtmlBody("https://www.pinkbike.com/buysell/list/?region=3&q=maverick&framesize=9,11,12,17,18,20,21,22&material=2"), modelPage: "maverick", },
		{ urlBodyPage: getHtmlBody("https://www.pinkbike.com/buysell/list/?region=3&q=sentinel&framesize=9,11,12,17,18,20,21,22"), modelPage: "sentinel", },
		{ urlBodyPage: getHtmlBody("https://www.pinkbike.com/buysell/list/?region=3&q=bronson&framesize=9,11,12,17,18,20,21,22&material=2"), modelPage: "bronson", },
		{ urlBodyPage: getHtmlBody("https://www.pinkbike.com/buysell/list/?region=3&q=roubion&framesize=9,11,12,17,18,20,21,22"), modelPage: "roubion",},
	}
	
	for _, bikePage := range bikesToSearch {
		bikePageArray := findTheBikes(bikePage.urlBodyPage, bikePage.modelPage)
		bikesArray = append(bikePageArray)
	}

	if(len(bikesArray) == 0 || hour == 4) {
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
