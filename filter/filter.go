package filter

import (
	"reflect"
	"strings"
	"sync"
)

type Trail struct {
	AccessName string
	RESTROOMS  string
	PICNIC     string
	FISHING    string
	Fee        string
	BikeRack   string
	BikeTrail  string
	DogTube    string
	Grills     string
	ADAtoilet  string
	ADAfishing string
	ADAcamping string
	ADApicnic  string
	ADAtrail   string
	ADAparking string
	ADAfacilit string
	HorseTrail string
	RecycleBin string
	DogCompost string
	THLeash    string
}

func matchesPartial(fieldValue, filterValue string) bool {
	return strings.Contains(strings.ToLower(strings.TrimSpace(fieldValue)), strings.ToLower(strings.TrimSpace(filterValue)))
}

func filterChunk(trails []Trail, filters map[string]string, resultChan chan []Trail, wg *sync.WaitGroup) {
	defer wg.Done()
	var filtered []Trail
	for _, trail := range trails {
		match := true
		for key, filterValue := range filters {
			if filterValue != "" {
				trailValue := reflect.ValueOf(trail).FieldByName(key)
				if !trailValue.IsValid() {
					match = false
					break
				}
				if !matchesPartial(trailValue.String(), filterValue) {
					match = false
					break
				}
			}
		}
		if match {
			filtered = append(filtered, trail)
		}
	}

	resultChan <- filtered
}

func FilterTrailsParallel(trails []Trail, filters map[string]string) []Trail {
	var wg sync.WaitGroup
	resultChan := make(chan []Trail, len(trails)/10)
	chunkSize := len(trails) / 4
	if chunkSize == 0 {
		chunkSize = 1
	}

	for i := 0; i < len(trails); i += chunkSize {
		end := i + chunkSize
		if end > len(trails) {
			end = len(trails)
		}
		wg.Add(1)
		go filterChunk(trails[i:end], filters, resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var finalResults []Trail
	for chunk := range resultChan {
		finalResults = append(finalResults, chunk...)
	}

	return finalResults
}
