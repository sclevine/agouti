package element

import (
	"fmt"
	"github.com/sclevine/agouti/core/internal/types"
	"strings"
)

type Element struct {
	ID      string
	Session session
}

type session interface {
	Execute(endpoint, method string, body interface{}, result ...interface{}) error
}

func (e *Element) GetID() string {
	return e.ID
}

func (e *Element) GetElements(selector types.Selector) ([]types.Element, error) {
	var results []struct{ Element string }

	if err := e.Session.Execute(e.url()+"/elements", "POST", selector, &results); err != nil {
		return nil, err
	}

	elements := []types.Element{}
	for _, result := range results {
		elements = append(elements, &Element{result.Element, e.Session})
	}

	return elements, nil
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

func (e *Element) GetCSS(property string) (string, error) {
	var value string
	if err := e.Session.Execute(fmt.Sprintf("%s/css/%s", e.url(), property), "GET", nil, &value); err != nil {
		return "", err
	}
	return value, nil
}

func (e *Element) Click() error {
	return e.Session.Execute(e.url()+"/click", "POST", nil, &struct{}{})
}

func (e *Element) Clear() error {
	return e.Session.Execute(e.url()+"/clear", "POST", nil, &struct{}{})
}

func (e *Element) Value(text string) error {
	splitText := strings.Split(text, "")
	request := struct {
		Value []string `json:"value"`
	}{splitText}
	return e.Session.Execute(e.url()+"/value", "POST", request, &struct{}{})
}

func (e *Element) IsSelected() (bool, error) {
	var selected bool
	if err := e.Session.Execute(e.url()+"/selected", "GET", nil, &selected); err != nil {
		return false, err
	}
	return selected, nil
}

func (e *Element) IsDisplayed() (bool, error) {
	var displayed bool
	if err := e.Session.Execute(e.url()+"/displayed", "GET", nil, &displayed); err != nil {
		return false, err
	}
	return displayed, nil
}

func (e *Element) IsEnabled() (bool, error) {
	var enabled bool
	if err := e.Session.Execute(e.url()+"/enabled", "GET", nil, &enabled); err != nil {
		return false, err
	}
	return enabled, nil
}

func (e *Element) Submit() error {
	return e.Session.Execute(e.url()+"/submit", "POST", nil, &struct{}{})
}

func (e *Element) url() string {
	return "element/" + e.ID
}

func (e *Element) IsEqualTo(other types.Element) (bool, error) {
	var equal bool
	if err := e.Session.Execute(e.url()+"/equals/"+other.GetID(), "GET", nil, &equal); err != nil {
		return false, err
	}
	return equal, nil
}
