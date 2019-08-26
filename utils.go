package gdriver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// WriteJSON - Write to json file
func WriteJSON(obj interface{}, filename string) {
	data, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

// ReadJSON - Read JSON
func ReadJSON(filename string) []map[string]interface{} {
	// Open our jsonFile
	jsonFile, err := os.Open(filename)

	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	return result
}
