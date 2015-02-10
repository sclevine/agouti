package api

import (
	"path"
	"strings"
)

type Element struct {
	ID      string
	Session *Session
}

func (e *Element) GetElement(selector Selector) (*Element, error) {
	var result struct{ Element string }

	if err := e.Session.sendElement(e.ID, "element", "POST", selector, &result); err != nil {
		return nil, err
	}

	return &Element{result.Element, e.Session}, nil
}

func (e *Element) GetElements(selector Selector) ([]*Element, error) {
	var results []struct{ Element string }

	if err := e.Session.sendElement(e.ID, "elements", "POST", selector, &results); err != nil {
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
	if err := e.Session.sendElement(e.ID, "text", "GET", nil, &text); err != nil {
		return "", err
	}
	return text, nil
}

func (e *Element) GetAttribute(attribute string) (string, error) {
	var value string
	if err := e.Session.sendElement(e.ID, path.Join("attribute", attribute), "GET", nil, &value); err != nil {
		return "", err
	}
	return value, nil
}

func (e *Element) GetCSS(property string) (string, error) {
	var value string
	if err := e.Session.sendElement(e.ID, path.Join("css", property), "GET", nil, &value); err != nil {
		return "", err
	}
	return value, nil
}

func (e *Element) Click() error {
	return e.Session.sendElement(e.ID, "click", "POST", nil)
}

func (e *Element) Clear() error {
	return e.Session.sendElement(e.ID, "clear", "POST", nil)
}

func (e *Element) Value(text string) error {
	splitText := strings.Split(text, "")
	request := struct {
		Value []string `json:"value"`
	}{splitText}
	return e.Session.sendElement(e.ID, "value", "POST", request)
}

func (e *Element) IsSelected() (bool, error) {
	var selected bool
	if err := e.Session.sendElement(e.ID, "selected", "GET", nil, &selected); err != nil {
		return false, err
	}
	return selected, nil
}

func (e *Element) IsDisplayed() (bool, error) {
	var displayed bool
	if err := e.Session.sendElement(e.ID, "displayed", "GET", nil, &displayed); err != nil {
		return false, err
	}
	return displayed, nil
}

func (e *Element) IsEnabled() (bool, error) {
	var enabled bool
	if err := e.Session.sendElement(e.ID, "enabled", "GET", nil, &enabled); err != nil {
		return false, err
	}
	return enabled, nil
}

func (e *Element) Submit() error {
	return e.Session.sendElement(e.ID, "submit", "POST", nil)
}

func (e *Element) IsEqualTo(other *Element) (bool, error) {
	var equal bool
	if err := e.Session.sendElement(e.ID, path.Join("equals", other.ID), "GET", nil, &equal); err != nil {
		return false, err
	}
	return equal, nil
}
