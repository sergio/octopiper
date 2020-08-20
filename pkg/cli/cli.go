package cli

import (
	"encoding/json"
	"fmt"
)

// WriteOutput is a func that writes json to stdout
func WriteOutput(jsondata interface{}) {
	result, err := json.MarshalIndent(jsondata, "", "   ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(result))
}
