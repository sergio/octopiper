package octopusapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// OctopusAPI :
type OctopusAPI struct {
	BaseURL string
	APIKey  string
}

// Get :
func (c OctopusAPI) Get(resource string) (string, error) {

	url := fmt.Sprintf("%s/api/%s", c.BaseURL, resource)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("X-Octopus-ApiKey", c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("error invoking octopus API: %w", err)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading octopus response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("expected 200 OK, but was '%s': %s", resp.Status, content)
	}

	return string(content), nil
}
