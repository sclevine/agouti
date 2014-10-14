package webdriver

import (
	"encoding/base64"
	"github.com/sclevine/agouti/core/internal/webdriver/element"
	"github.com/sclevine/agouti/core/internal/webdriver/window"
)

type Driver struct {
	Session executable
}

type executable interface {
	Execute(endpoint, method string, body, result interface{}) error
}

type Element interface {
	GetID() string
	GetText() (string, error)
	GetAttribute(attribute string) (string, error)
	GetCSS(property string) (string, error)
	IsSelected() (bool, error)
	IsDisplayed() (bool, error)
	Click() error
	Clear() error
	Value(text string) error
	Submit() error
}

type Window interface {
	SetSize(height, width int) error
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

func (d *Driver) GetElements(selector string) ([]Element, error) {
	request := struct {
		Using string `json:"using"`
		Value string `json:"value"`
	}{"css selector", selector}

	var results []struct{ Element string }

	if err := d.Session.Execute("elements", "POST", request, &results); err != nil {
		return nil, err
	}

	elements := []Element{}
	for _, result := range results {
		elements = append(elements, &element.Element{result.Element, d.Session})
	}

	return elements, nil
}

func (d *Driver) GetWindow() (Window, error) {
	var windowID string
	if err := d.Session.Execute("window_handle", "GET", nil, &windowID); err != nil {
		return nil, err
	}
	return &window.Window{windowID, d.Session}, nil
}

func (d *Driver) SetCookie(cookie *Cookie) error {
	request := struct {
		Cookie *Cookie `json:"cookie"`
	}{cookie}

	return d.Session.Execute("cookie", "POST", request, &struct{}{})
}

func (d *Driver) DeleteCookie(cookieName string) error {
	return d.Session.Execute("cookie/"+cookieName, "DELETE", nil, &struct{}{})
}

func (d *Driver) DeleteCookies() error {
	return d.Session.Execute("cookie", "DELETE", nil, &struct{}{})
}

func (d *Driver) GetScreenshot() ([]byte, error) {
	var base64Image string

	if err := d.Session.Execute("screenshot", "GET", nil, &base64Image); err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(base64Image)
}

func (d *Driver) GetURL() (string, error) {
	var url string
	if err := d.Session.Execute("url", "GET", nil, &url); err != nil {
		return "", err
	}

	return url, nil
}

func (d *Driver) SetURL(url string) error {
	request := struct {
		URL string `json:"url"`
	}{url}

	return d.Session.Execute("url", "POST", request, &struct{}{})
}

func (d *Driver) GetTitle() (string, error) {
	var title string
	if err := d.Session.Execute("title", "GET", nil, &title); err != nil {
		return "", err
	}

	return title, nil
}

func (d *Driver) DoubleClick() error {
	return d.Session.Execute("doubleclick", "POST", nil, &struct{}{})
}

func (d *Driver) MoveTo(element Element, point Point) error {
	request := map[string]interface{}{}

	if element != nil {
		request["element"] = element.GetID()
	}

	if point != nil {
		if xoffset, present := point.x(); present {
			request["xoffset"] = xoffset
		}

		if yoffset, present := point.y(); present {
			request["yoffset"] = yoffset
		}
	}

	return d.Session.Execute("moveto", "POST", request, &struct{}{})
}
