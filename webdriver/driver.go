package webdriver

import (
	"fmt"
	"github.com/sclevine/agouti/webdriver/element"
)

type Executable interface {
	Execute(endpoint, method string, body, result interface{}) error
}

type Driver struct {
	Session Executable
}

type Element interface {
	GetText() (string, error)
}

type Cookie struct {
	Name     string      `json:"name"`
	Value    interface{} `json:"value"`
	Path     string      `json:"path"`
	Domain   string      `json:"domain"`
	Secure   bool        `json:"secure"`
	HTTPOnly bool        `json:"httpOnly"`
	Expiry   int64       `json:"expiry"`
}

func (d *Driver) Navigate(url string) error {
	request := struct {
		URL string `json:"url"`
	}{url}
	if err := d.Session.Execute("url", "POST", request, &struct{}{}); err != nil {
		return fmt.Errorf("failed to navigate: %s", err)
	}
	return nil
}

func (d *Driver) GetElements(selector string) ([]Element, error) {
	request := struct {
		Using string `json:"using"`
		Value string `json:"value"`
	}{"css selector", selector}

	var results []struct{ Element string }

	if err := d.Session.Execute("elements", "POST", request, &results); err != nil {
		return nil, fmt.Errorf("failed to get elements with selector '%s': %s", selector, err)
	}

	elements := []Element{}
	for _, result := range results {
		elements = append(elements, &element.Element{result.Element, d.Session})
	}

	return elements, nil
}

func (d *Driver) SetCookie(cookie *Cookie) error {
	request := struct {
		Cookie *Cookie `json:"cookie"`
	}{cookie}

	if err := d.Session.Execute("cookie", "POST", request, &struct{}{}); err != nil {
		return fmt.Errorf("failed to add cookie: %s", err)
	}

	return nil
}
