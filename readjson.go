package gdriver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func ReadJSON(filename string) []map[string]string {
	// Open our jsonFile
	jsonFile, err := os.Open(filename)

	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []map[string]string
	json.Unmarshal([]byte(byteValue), &result)

	return result
}
