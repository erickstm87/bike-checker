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
    // comment
	theDate := strconv.Itoa(year) + "-" + month.String() + "-" + strconv.Itoa(day) + " " + strconv.Itoa(hour) + ":" + sMinute

	fmt.Println("starting job at: ", theDate)
	var bikesArray []AvailableBike
	searches := []struct {
		url   string
		model string
	}{
		{url: "https://www.pinkbike.com/buysell/list/?region=3&q=hightower&framesize=9,11,12,17,18,20,21,22", model: "hightower"},
		{url: "https://www.pinkbike.com/buysell/list/?region=3&q=maverick&framesize=9,11,12,17,18,20,21,22&material=2", model: "maverick"},
		{url: "https://www.pinkbike.com/buysell/list/?region=3&q=sentinel&framesize=9,11,12,17,18,20,21,22", model: "sentinel"},
		{url: "https://www.pinkbike.com/buysell/list/?region=3&q=bronson&framesize=9,11,12,17,18,20,21,22&material=2", model: "bronson"},
		{url: "https://www.pinkbike.com/buysell/list/?region=3&q=roubion&framesize=9,11,12,17,18,20,21,22", model: "roubion"},
	}

	for _, search := range searches {
		doc, err := getHtmlBody(search.url)
		if err != nil {
			fmt.Printf("Error fetching %s: %v\n", search.model, err)
			continue
		}
		foundBikes := findTheBikes(doc, search.model)
		bikesArray = append(bikesArray, foundBikes...)
	}

	if len(bikesArray) == 0 || hour == 4 {
		fmt.Println("no entries were found")
		return
	}
	readDB(bikesArray)
	// seedDB(bikesAvailable)
}

func errorHandler(err error) {
	fmt.Println("there was an error: ", err)
}

func getHtmlBody(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status %d for URL %s", resp.StatusCode, url)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML from %s: %w", url, err)
	}
	return doc, nil
}

func findTheBikes(bikeEntries *goquery.Document, modelName string) []AvailableBike {
	var bikes []AvailableBike
	bikeEntries.Find(".uImage").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Children().Attr("href")
		if !exists || link == "" {
			return
		}
		
		// Skip known bad/test entries
		if link == "https://www.pinkbike.com/buysell/3029717/" {
			return
		}

		bikes = append(bikes, AvailableBike{
			model: modelName,
			link:  link,
		})
	})
	return bikes
}
