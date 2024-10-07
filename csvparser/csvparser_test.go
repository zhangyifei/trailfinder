package csvparser

import (
	"reflect"
	"strings"
	"testing"
	"trailfinder/filter"
)

func collectResults(trailChan <-chan []filter.Trail, doneChan <-chan error) ([]filter.Trail, error) {
	var allTrails []filter.Trail
	for {
		select {
		case trails, ok := <-trailChan:
			if !ok {
				trailChan = nil // Mark trailChan as nil to avoid further reads
			} else {
				allTrails = append(allTrails, trails...)
			}
		case err := <-doneChan:
			if err != nil {
				return allTrails, err
			}
			doneChan = nil // Mark doneChan as nil to avoid further reads
		}

		// Exit the loop when both channels are nil
		if trailChan == nil && doneChan == nil {
			break
		}
	}
	return allTrails, nil
}

func TestParseCSVConcurrently(t *testing.T) {
	// Sample CSV content
	csvContent := `AccessName,RESTROOMS,PICNIC
Trail1,Yes,No
Trail2,No,Yes
`

	// Use a strings.Reader to simulate a file
	reader := strings.NewReader(csvContent)

	// Mock channels
	trailChan := make(chan []filter.Trail)
	doneChan := make(chan error)

	// Call the function with chunk size 2
	go ParseCSVConcurrently(reader, 2, trailChan, doneChan)

	// Collect the results
	trails, err := collectResults(trailChan, doneChan)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the parsed result
	expectedTrails := []filter.Trail{
		{AccessName: "Trail1", RESTROOMS: "Yes", PICNIC: "No"},
		{AccessName: "Trail2", RESTROOMS: "No", PICNIC: "Yes"},
	}

	if !reflect.DeepEqual(trails, expectedTrails) {
		t.Errorf("Expected %v, got %v", expectedTrails, trails)
	}
}
