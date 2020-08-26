package octopusclient

import (
	"encoding/json"
	"fmt"
	"log"
	"octopiper/pkg/octopusapi"
	"octopiper/pkg/sqlitekvs"
)

// OctopusClient :
type OctopusClient struct {
	LocalCache    *sqlitekvs.KVSCache
	OctopusServer *octopusapi.OctopusAPI
}

// GetJSON :
func (c *OctopusClient) GetJSON(resource string) (interface{}, error) {
	content, err := c.Get(resource)
	if err != nil {
		return "", err
	}

	var jsondata interface{}
	err = json.Unmarshal([]byte(content), &jsondata)
	if err != nil {
		return "", fmt.Errorf("unable to deserialize as json: %w", err)
	}

	return jsondata, nil
}

// Get :
func (c *OctopusClient) Get(resource string) (string, error) {

	cacheKey := resource

	content, found, err := c.LocalCache.Get(cacheKey)
	if err != nil {
		return "", fmt.Errorf("could not query local cache for '%s': %w", cacheKey, err)
	}

	if found {
		return content, nil
	}

	content, err = c.OctopusServer.Get(resource)
	if err != nil {
		return "", fmt.Errorf("cound not download content from octopus for '%s': %w", resource, err)
	}

	err = c.LocalCache.Set(cacheKey, content)
	if err != nil {
		log.Printf("could not cache content: %v", err)
	}
	return content, nil
}
