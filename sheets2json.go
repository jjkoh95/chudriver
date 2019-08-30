package gdriver

import (
	"encoding/csv"
	"io"
	"net/http"
)

// ExportCSVSheet - The easiest way to make it work - Publish as CSV
func ExportCSVSheet(spreadsheetURL, jsonFilename string) {
	obj := ReadCSVSheet(spreadsheetURL)
	WriteJSON(obj, jsonFilename)
}

// ReadCSVSheet - Read CSV sheet - return object
func ReadCSVSheet(spreadsheetURL string) []map[string]interface{} {
	resp, err := http.Get(spreadsheetURL)
	if err != nil {
		panic("Error reading spreadsheet")
	}
	reader := csv.NewReader(io.Reader(resp.Body))
	obj, err := reader.ReadAll()
	if err != nil {
		panic("Error reading content of spreadsheet")
	}
	return convertRowsToJSON(obj)
}
