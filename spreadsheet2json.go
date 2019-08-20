package gdriver

import (
	"encoding/csv"
	"io"
	"net/http"
)

// CSVResponse2JSON - csv response to JSON
func CSVResponse2JSON(resp *http.Response) interface{} {
	reader := csv.NewReader(io.Reader(resp.Body))
	val, err := reader.ReadAll()

	if err != nil {
		panic("Error reading response")
	}

	obj := make([]map[string]string, len(val)-1)
	keys := val[0]
	for row := 1; row < len(val); row++ {
		obj[row-1] = make(map[string]string) // need to explicitly declare map again
		for col := 0; col < len(keys); col++ {
			obj[row-1][keys[col]] = val[row][col]
		}
	}

	return obj
}
