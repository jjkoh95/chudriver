package gdriver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

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
