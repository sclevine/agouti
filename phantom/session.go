package phantom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Session struct {
	URL string
}

func (s *Session) Execute(endpoint, method string, body, result interface{}) error {
	client := &http.Client{}

	var bodyReader io.Reader
	if body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("invalid request body: %s", err)
		}
		bodyReader = bytes.NewReader(bodyJSON)
	}

	request, err := http.NewRequest(method, s.URL+"/"+endpoint, bodyReader)
	if err != nil {
		return fmt.Errorf("invalid request: %s", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("request failed: %s", err)
	}

	responseBody, _ := ioutil.ReadAll(response.Body)
	if err := json.Unmarshal(responseBody, result); err != nil {
		return fmt.Errorf("invalid response body: %s", err)
	}

	return nil
}
