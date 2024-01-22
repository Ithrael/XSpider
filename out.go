package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func createCsvWriter() (*csv.Writer, *os.File, error) {
	csvFile, err := os.OpenFile(out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	writer := csv.NewWriter(bufio.NewWriter(csvFile))
	return writer, csvFile, nil
}

func WriteDetailsToCSV(detailsCh chan *PageDetail) {
	writer, f, err := createCsvWriter()
	if err != nil {
		log.Fatalf("Failed to create CSV writer: %v", err)
	}
	defer f.Close()
	defer writer.Flush()

	for detail := range detailsCh {
		err := writer.Write([]string{
			detail.Url,
			detail.Title,
			fmt.Sprint(detail.ResponseCode),
			detail.Fingerprint,
		})
		if err != nil {
			log.Printf("Failed to write data to CSV: %v", err)
		}
	}
}
