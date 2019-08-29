package gdriver

import (
	"strconv"
	"strings"

	"google.golang.org/api/docs/v1"
)

// ColType - store dtype and parseFunc
type ColType struct {
	dtype     int
	parseFunc ParseFuncType
}

// ParseFuncType - wrap parsing function type
type ParseFuncType = func(s string) interface{}

func parseFloat(s string) interface{} {
	parsedVal, _ := strconv.ParseFloat(s, 64)
	return parsedVal
}

func parseInt(s string) interface{} {
	parsedVal, _ := strconv.Atoi(s)
	return parsedVal
}

func parseBool(s string) interface{} {
	return strings.ToLower(s) == "true"
}

func parseString(s string) interface{} {
	return s
}

// parseInterface2String - this makes [][]interface{} much easier to work with
func parseInterface2String(val [][]interface{}) [][]string {
	resp := make([][]string, len(val))
	for i := 0; i < len(val); i++ {
		resp[i] = make([]string, len(val[i]))
		for j := 0; j < len(val[i]); j++ {
			resp[i][j] = val[i][j].(string)
		}
	}
	return resp
}

// convertRowsToJSON - convert [][]string to []map[string] (json object)
func convertRowsToJSON(val [][]string) []map[string]interface{} {
	keys := val[0] // header is always key
	colTypes := make(map[string]ColType)
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
			colTypes[keys[col]] = ColType{
				dtype:     DSTRING,
				parseFunc: parseString,
			}
		} else if hasFloat {
			colTypes[keys[col]] = ColType{
				dtype:     DFLOAT,
				parseFunc: parseFloat,
			}
		} else if hasInt {
			colTypes[keys[col]] = ColType{
				dtype:     DINT,
				parseFunc: parseInt,
			}
		} else if hasBool {
			colTypes[keys[col]] = ColType{
				dtype:     DBOOL,
				parseFunc: parseBool,
			}
		}
	}

	obj := make([]map[string]interface{}, len(val)-1) // this is the object we want to return
	for row := 1; row < len(val); row++ {
		obj[row-1] = make(map[string]interface{})
		for col := 0; col < len(keys); col++ {
			if val[row][col] == "" || val[row][col] == "nan" {
				obj[row-1][keys[col]] = nil
				continue
			}
			parsedVal := colTypes[keys[col]].parseFunc(val[row][col])
			obj[row-1][keys[col]] = parsedVal
		}
	}

	return obj
}

// ReadStructuralElements - read docs structural elements to string
func readStructuralElements(struturalElements []*docs.StructuralElement) string {
	var sb strings.Builder
	for _, structuralElement := range struturalElements {
		if structuralElement.Paragraph == nil {
			continue
		}
		for _, paragraphElement := range structuralElement.Paragraph.Elements {
			textRun := paragraphElement.TextRun
			sb.WriteString(textRun.Content)
		}
	}
	return sb.String()
}
