package element

import "fmt"

type Executable interface {
	Execute(endpoint, method string, body, result interface{}) error
}

type Element struct {
	ID      string
	Session Executable
}

func (e *Element) GetText() (string, error) {
	var text string
	if err := e.Session.Execute("element/"+e.ID+"/text", "GET", nil, &text); err != nil {
		return "", fmt.Errorf("failed to retrieve text: %s", err)
	}
	return text, nil
}
