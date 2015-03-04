package api

import (
	"path"
	"strings"
)

type Element struct {
	ID      string
	Session *Session
}

func (e *Element) Send(endpoint, method string, body, result interface{}) error {
	return e.Session.Send(path.Join("element", e.ID, endpoint), method, body, result)
}

func (e *Element) GetElement(selector Selector) (*Element, error) {
	var result struct{ Element string }

	if err := e.Send("element", "POST", selector, &result); err != nil {
		return nil, err
	}

	return &Element{result.Element, e.Session}, nil
}

func (e *Element) GetElements(selector Selector) ([]*Element, error) {
	var results []struct{ Element string }

	if err := e.Send("elements", "POST", selector, &results); err != nil {
		return nil, err
	}

	elements := []*Element{}
	for _, result := range results {
		elements = append(elements, &Element{result.Element, e.Session})
	}

	return elements, nil
}

func (e *Element) GetText() (string, error) {
	var text string
	if err := e.Send("text", "GET", nil, &text); err != nil {
		return "", err
	}
	return text, nil
}

func (e *Element) GetName() (string, error) {
	var name string
	if err := e.Send("name", "GET", nil, &name); err != nil {
		return "", err
	}
	return name, nil
}

func (e *Element) GetAttribute(attribute string) (string, error) {
	var value string
	if err := e.Send(path.Join("attribute", attribute), "GET", nil, &value); err != nil {
		return "", err
	}
	return value, nil
}

func (e *Element) GetCSS(property string) (string, error) {
	var value string
	if err := e.Send(path.Join("css", property), "GET", nil, &value); err != nil {
		return "", err
	}
	return value, nil
}

func (e *Element) Click() error {
	return e.Send("click", "POST", nil, nil)
}

func (e *Element) Clear() error {
	return e.Send("clear", "POST", nil, nil)
}

func (e *Element) Value(text string) error {
	splitText := strings.Split(text, "")
	request := struct {
		Value []string `json:"value"`
	}{splitText}
	return e.Send("value", "POST", request, nil)
}

func (e *Element) IsSelected() (bool, error) {
	var selected bool
	if err := e.Send("selected", "GET", nil, &selected); err != nil {
		return false, err
	}
	return selected, nil
}

func (e *Element) IsDisplayed() (bool, error) {
	var displayed bool
	if err := e.Send("displayed", "GET", nil, &displayed); err != nil {
		return false, err
	}
	return displayed, nil
}

func (e *Element) IsEnabled() (bool, error) {
	var enabled bool
	if err := e.Send("enabled", "GET", nil, &enabled); err != nil {
		return false, err
	}
	return enabled, nil
}

func (e *Element) Submit() error {
	return e.Send("submit", "POST", nil, nil)
}

func (e *Element) IsEqualTo(other *Element) (bool, error) {
	var equal bool
	if err := e.Send(path.Join("equals", other.ID), "GET", nil, &equal); err != nil {
		return false, err
	}
	return equal, nil
}
