package csvparser

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"trailfinder/filter"
)

func ParseCSVConcurrently(reader io.Reader, chunkSize int, trailChan chan<- []filter.Trail, doneChan chan<- error) {
	defer close(trailChan)

	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1

	header, err := csvReader.Read()
	if err != nil {
		doneChan <- fmt.Errorf("error reading CSV header: %w", err)
		return
	}

	fieldIndices := make(map[string]int)
	for i, field := range header {
		fieldIndices[field] = i
	}

	var trails []filter.Trail
	for {
		record, err := csvReader.Read()
		if err == csv.ErrFieldCount {
			continue
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			doneChan <- fmt.Errorf("error reading CSV record: %w", err)
			return
		}

		trail := filter.Trail{}
		for key, index := range fieldIndices {
			if index < len(record) {
				trailValue := reflect.ValueOf(&trail).Elem().FieldByName(key)
				if trailValue.IsValid() && trailValue.CanSet() {
					trailValue.SetString(record[index])
				}
			}
		}
		trails = append(trails, trail)
		if len(trails) >= chunkSize {
			trailChan <- trails
			trails = nil
		}
	}

	if len(trails) > 0 {
		trailChan <- trails
	}
	doneChan <- nil
}
