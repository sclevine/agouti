package webdriver

import "fmt"

type Executable interface {
	Execute(endpoint, method string, body, result interface{}) error
}

type Driver struct {
	Session Executable
}

type Element struct {
	ID string
}

func (d *Driver) Navigate(url string) error {
	request := struct{URL string `json:"url"`}{url}
	if err := d.Session.Execute("url", "POST", request, &struct{}{}); err != nil {
		return fmt.Errorf("failed to navigate: %s", err)
	}
	return nil
}

func(d *Driver) GetElements(selector string) []*Element {
	request := struct{
		Using string `json:"using"`
		Value string `json:"value"`
	}{ "css selector", selector }

	result := struct{
		
	}{}

	if err := d.Session.Execute("elements", "POST", request, ); err != nil {
		return fmt.Errorf("failed to navigate: %s", err)
	}
	return nil
	return nil
}

func (e *Element) GetText() string {
	return ""
}
