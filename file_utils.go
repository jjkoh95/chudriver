package gdriver

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// WriteJSON - Write to json file
func WriteJSON(obj interface{}, filename string) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// ReadJSON - Read JSON
func ReadJSON(filename string) ([]map[string]interface{}, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(filename)

	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil
}

// WriteText - Write raw text
func WriteText(text, filename string) error {
	return ioutil.WriteFile(filename, []byte(text), 0644)
}
