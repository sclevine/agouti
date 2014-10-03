package webdriver

import "fmt"

type Executable interface {
	Execute(endpoint, method string, body, result interface{}) error
}

type Driver struct {
	Session Executable
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

func (d *Driver) GetElements(selector string) ([]*Element, error) {
	request := struct {
		Using string `json:"using"`
		Value string `json:"value"`
	}{"css selector", selector}

	var results []struct{ Element string }

	if err := d.Session.Execute("elements", "POST", request, &results); err != nil {
		return nil, fmt.Errorf("failed to get elements with selector '%s': %s", selector, err)
	}

	elements := []*Element{}
	for _, result := range results {
		elements = append(elements, &Element{result.Element, d.Session})
	}

	return elements, nil
}
