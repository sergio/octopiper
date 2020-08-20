package octopusapi

import (
	"fmt"
	"octopiper/pkg/rest"
)

// Server is a struct
type Server struct {
	BaseURL string
	APIKey  string
}

// GetJSON is a func
func (octopusServer *Server) GetJSON(resource string) (interface{}, error) {
	endpoint := fmt.Sprintf("%s/api/%s?apikey=%s", octopusServer.BaseURL, resource, octopusServer.APIKey)
	return rest.GetJSON(endpoint)
}
