package main

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestFindTheBikes(t *testing.T) {
	// Test HTML containing bike listings
	html := `
		<html>
			<body>
				<div class="uImage">
					<a href="https://www.pinkbike.com/buysell/1234/">Good Link</a>
				</div>
				<div class="uImage">
					<a href="https://www.pinkbike.com/buysell/3029717/">Excluded Link</a>
				</div>
				<div class="uImage">
					<a href="https://www.pinkbike.com/buysell/5678/">Another Good Link</a>
				</div>
				<div class="uImage">
					<span>Invalid - No Link</span>
				</div>
			</body>
		</html>
	`

	// Parse the test HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("Failed to parse test HTML: %v", err)
	}

	// Test the findTheBikes function
	bikes := findTheBikes(doc, "test-model")

	// Should find exactly 2 valid bikes (excluding the known bad link and invalid entry)
	if len(bikes) != 2 {
		t.Errorf("Expected 2 bikes, got %d", len(bikes))
	}

	// Check the first bike
	if bikes[0].link != "https://www.pinkbike.com/buysell/1234/" {
		t.Errorf("Expected first link to be https://www.pinkbike.com/buysell/1234/, got %s", bikes[0].link)
	}
	if bikes[0].model != "test-model" {
		t.Errorf("Expected model to be test-model, got %s", bikes[0].model)
	}

	// Check the second bike
	if bikes[1].link != "https://www.pinkbike.com/buysell/5678/" {
		t.Errorf("Expected second link to be https://www.pinkbike.com/buysell/5678/, got %s", bikes[1].link)
	}
}

func TestGetHtmlBody(t *testing.T) {
	// Test with invalid URL
	_, err := getHtmlBody("http://invalid.url.that.should.fail")
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}

	// Test with valid URL but might be down/unavailable
	doc, err := getHtmlBody("https://www.pinkbike.com")
	if err != nil {
		t.Logf("Note: Connection test failed (this might be normal if no internet): %v", err)
		return
	}

	// If we got here, check we got a valid document
	if doc == nil {
		t.Error("Expected non-nil document for valid URL")
	}
}