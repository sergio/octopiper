package cmd

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {

	globalconfig.OctopusServer.BaseURL = "https://octopus.reachcore.net"
	globalconfig.OctopusServer.APIKey = "API-JDIEZATYHNHAQFXVJJTLICHBLQ"
	globalconfig.LocalCache.FilePath = "test.db"
	globalconfig.LocalCache.TTLMinutes = 1

	results, err := find("Source")
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, results)

	/*
		for _, r := range results {
			fmt.Printf("%s: %s %s\n", r.Name, r.Value, r.Scope)
		}
	*/

	jsonText, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(jsonText))
}
