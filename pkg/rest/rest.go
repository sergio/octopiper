package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// GetJSON returns json from a url endpoint
func GetJSON(endpoint string) (interface{}, error) {

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jsonText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data interface{}
	err = json.Unmarshal(jsonText, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
