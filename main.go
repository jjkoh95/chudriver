package gdriver

import (
	"net/http"
)

func SpreadsheetCSVToJSON(spreadsheetURL string, jsonFilename string) {
	resp, _ := http.Get(spreadsheetURL)
	obj := CSVResponse2JSON(resp)
	WriteJSON(obj, jsonFilename)
}
