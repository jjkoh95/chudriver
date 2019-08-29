package gdriver

import (
	"encoding/csv"
	"io"
	"net/http"
)

// SpreadsheetCSVToJSON - This is used for public csv spreadsheet
func SpreadsheetCSVToJSON(spreadsheetURL string, jsonFilename string) {
	resp, _ := http.Get(spreadsheetURL)
	obj := CSVResponse2JSON(resp)
	WriteJSON(obj, jsonFilename)
}

// CSVResponse2JSON - csv response to JSON
func CSVResponse2JSON(resp *http.Response) interface{} {
	reader := csv.NewReader(io.Reader(resp.Body))
	val, err := reader.ReadAll()

	if err != nil {
		panic("Error reading response")
	}

	return convertRowsToJSON(val)
}
