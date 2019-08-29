package gdriver

import (
	"encoding/csv"
	"io"
	"net/http"
)

// SpreadsheetCSVToJSON - The easiest way to make it work - Publish as CSV
func SpreadsheetCSVToJSON(spreadsheetURL string, jsonFilename string) {
	resp, err := http.Get(spreadsheetURL)
	if err != nil {
		panic("Error reading spreadsheet")
	}
	reader := csv.NewReader(io.Reader(resp.Body))
	obj, err := reader.ReadAll()
	if err != nil {
		panic("Error reading content of spreadsheet")
	}
	WriteJSON(obj, jsonFilename)
}
