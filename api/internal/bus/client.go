package bus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	SessionURL string
}

func (c *Client) Send(endpoint, method string, body interface{}, result ...interface{}) error {
	client := &http.Client{}

	var bodyReader io.Reader
	if body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("invalid request body: %s", err)
		}
		bodyReader = bytes.NewReader(bodyJSON)
	}

	requestURL := strings.TrimSuffix(c.SessionURL+"/"+endpoint, "/")
	request, err := http.NewRequest(method, requestURL, bodyReader)
	if err != nil {
		return fmt.Errorf("invalid request: %s", err)
	}

	if body != nil {
		request.Header.Add("Content-Type", "application/json")
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("request failed: %s", err)
	}

	responseBody, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode < 200 || response.StatusCode > 299 {
		var errBody struct{ Value struct{ Message string } }
		if err := json.Unmarshal(responseBody, &errBody); err != nil {
			return fmt.Errorf("request unsuccessful: error unreadable")
		}

		var errMessage struct{ ErrorMessage string }
		if err := json.Unmarshal([]byte(errBody.Value.Message), &errMessage); err != nil {
			return fmt.Errorf("request unsuccessful: error message unreadable")
		}

		return fmt.Errorf("request unsuccessful: %s", errMessage.ErrorMessage)
	}

	if len(result) > 0 {
		bodyValue := struct{ Value interface{} }{result[0]}
		if err := json.Unmarshal(responseBody, &bodyValue); err != nil {
			return fmt.Errorf("failed to parse response value: %s", err)
		}
	}

	return nil
}

func Connect(url string, capabilities map[string]interface{}) (*Client, error) {
	if capabilities == nil {
		capabilities = map[string]interface{}{}
	}
	desiredCapabilities := struct {
		DesiredCapabilities map[string]interface{} `json:"desiredCapabilities"`
	}{capabilities}

	capabiltiesJSON, err := json.Marshal(desiredCapabilities)
	if err != nil {
		return nil, err
	}

	postBody := strings.NewReader(string(capabiltiesJSON))
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/session", url), postBody)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	var sessionResponse struct{ SessionID string }

	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &sessionResponse)

	if sessionResponse.SessionID == "" {
		return nil, errors.New("failed to retrieve a session ID")
	}

	sessionURL := fmt.Sprintf("%s/session/%s", url, sessionResponse.SessionID)
	return &Client{sessionURL}, nil
}
