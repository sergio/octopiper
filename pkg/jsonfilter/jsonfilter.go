package jsonfilter

import "github.com/jmespath/go-jmespath"

// Query is a function that takes a json structure and an optional jmespath query
func Query(query string, jsondata interface{}) (interface{}, error) {
	if query == "" {
		return jsondata, nil
	}
	jsondata, err := jmespath.Search(query, jsondata)
	if err != nil {
		return nil, err
	}
	return jsondata, nil
}
