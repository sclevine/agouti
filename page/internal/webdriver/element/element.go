package element

import "fmt"

type Element struct {
	ID      string
	Session executable
}

type executable interface {
	Execute(endpoint, method string, body, result interface{}) error
}

func (e *Element) GetText() (string, error) {
	var text string
	if err := e.Session.Execute(e.url()+"/text", "GET", nil, &text); err != nil {
		return "", err
	}
	return text, nil
}

func (e *Element) GetAttribute(attribute string) (string, error) {
	var value string
	if err := e.Session.Execute(fmt.Sprintf("%s/attribute/%s", e.url(), attribute), "GET", nil, &value); err != nil {
		return "", err
	}
	return value, nil
}

func (e *Element) Click() error {
	return e.Session.Execute(e.url()+"/click", "POST", nil, &struct{}{})
}

func (e *Element) url() string {
	return "element/" + e.ID
}
