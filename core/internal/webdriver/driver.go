package webdriver

import (
	"encoding/base64"
	"github.com/sclevine/agouti/core/internal/types"
	"github.com/sclevine/agouti/core/internal/webdriver/element"
	"github.com/sclevine/agouti/core/internal/webdriver/window"
)

type Driver struct {
	Session executable
}

type executable interface {
	Execute(endpoint, method string, body, result interface{}) error
}

func (d *Driver) GetElements(selector types.Selector) ([]types.Element, error) {
	var results []struct{ Element string }

	if err := d.Session.Execute("elements", "POST", selector, &results); err != nil {
		return nil, err
	}

	elements := []types.Element{}
	for _, result := range results {
		elements = append(elements, &element.Element{ID: result.Element, Session: d.Session})
	}

	return elements, nil
}

func (d *Driver) GetWindow() (types.Window, error) {
	var windowID string
	if err := d.Session.Execute("window_handle", "GET", nil, &windowID); err != nil {
		return nil, err
	}
	return &window.Window{ID: windowID, Session: d.Session}, nil
}

func (d *Driver) SetCookie(cookie *types.Cookie) error {
	request := struct {
		Cookie *types.Cookie `json:"cookie"`
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

func (d *Driver) MoveTo(element types.Element, point types.Point) error {
	request := map[string]interface{}{}

	if element != nil {
		request["element"] = element.GetID()
	}

	if point != nil {
		if xoffset, present := point.X(); present {
			request["xoffset"] = xoffset
		}

		if yoffset, present := point.Y(); present {
			request["yoffset"] = yoffset
		}
	}

	return d.Session.Execute("moveto", "POST", request, &struct{}{})
}

func (d *Driver) Execute(body string, arguments []interface{}, result interface{}) error {
	request := struct {
		Script string        `json:"script"`
		Args   []interface{} `json:"args"`
	}{body, arguments}

	if err := d.Session.Execute("execute", "POST", request, result); err != nil {
		return err
	}

	return nil
}

func (d *Driver) Forward() error {
	return d.Session.Execute("forward", "POST", nil, &struct{}{})
}

func (d *Driver) Back() error {
	return d.Session.Execute("back", "POST", nil, &struct{}{})
}

func (d *Driver) Refresh() error {
	return d.Session.Execute("refresh", "POST", nil, &struct{}{})
}
