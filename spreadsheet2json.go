package gdriver

import (
	"encoding/csv"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// CSVResponse2JSON - csv response to JSON
func CSVResponse2JSON(resp *http.Response) interface{} {
	reader := csv.NewReader(io.Reader(resp.Body))
	val, err := reader.ReadAll()

	if err != nil {
		panic("Error reading response")
	}

	keys := val[0] // header is always key
	colTypes := make(map[string]int)
	for col := 0; col < len(keys); col++ {
		var hasString, hasFloat, hasInt, hasBool bool
		for row := 1; row < len(val); row++ {
			if (hasBool && hasInt) || (hasBool && hasFloat) {
				break
			}
			if val[row][col] == "" || strings.ToLower(val[row][col]) == "nan" {
				continue
			}
			if _, err := strconv.Atoi(val[row][col]); err == nil {
				hasInt = true
				continue
			}
			if _, err := strconv.ParseFloat(val[row][col], 64); err == nil {
				hasFloat = true
				continue
			}
			if strings.ToLower(val[row][col]) == "true" || strings.ToLower(val[row][col]) == "false" {
				hasBool = true
				continue
			}
			hasString = true
			break
		}
		if hasString || (hasBool && hasInt) || (hasBool && hasFloat) {
			colTypes[keys[col]] = DSTRING
		} else if hasFloat {
			colTypes[keys[col]] = DFLOAT
		} else if hasInt {
			colTypes[keys[col]] = DINT
		} else if hasBool {
			colTypes[keys[col]] = DBOOL
		}
	}

	obj := make([]map[string]interface{}, len(val)-1)
	for row := 1; row < len(val); row++ {
		obj[row-1] = make(map[string]interface{}) // need to explicitly declare map again
		for col := 0; col < len(keys); col++ {
			if val[row][col] == "" || val[row][col] == "nan" {
				obj[row-1][keys[col]] = nil
			}
			var parsedVal interface{}
			if colTypes[keys[col]] == DFLOAT {
				parsedVal, _ = strconv.ParseFloat(val[row][col], 64)
			} else if colTypes[keys[col]] == DINT {
				parsedVal, _ = strconv.Atoi(val[row][col])
			} else if colTypes[keys[col]] == DBOOL {
				parsedVal = strings.ToLower(val[row][col]) == "true"
			} else {
				parsedVal = val[row][col]
			}
			obj[row-1][keys[col]] = parsedVal
		}
	}

	return obj
}
