package api

import (
	"encoding/base64"
	"github.com/sclevine/agouti/core/internal/api/element"
	"github.com/sclevine/agouti/core/internal/api/window"
	"github.com/sclevine/agouti/core/internal/types"
)

type Client struct {
	Session session
}

type session interface {
	Execute(endpoint, method string, body interface{}, result ...interface{}) error
}

func (c *Client) DeleteSession() error {
	return c.Session.Execute("", "DELETE", nil)
}

func (c *Client) GetElements(selector types.Selector) ([]types.Element, error) {
	var results []struct{ Element string }

	if err := c.Session.Execute("elements", "POST", selector, &results); err != nil {
		return nil, err
	}

	elements := []types.Element{}
	for _, result := range results {
		elements = append(elements, &element.Element{ID: result.Element, Session: c.Session})
	}

	return elements, nil
}

func (c *Client) GetWindow() (types.Window, error) {
	var windowID string
	if err := c.Session.Execute("window_handle", "GET", nil, &windowID); err != nil {
		return nil, err
	}
	return &window.Window{ID: windowID, Session: c.Session}, nil
}

func (c *Client) SetCookie(cookie *types.Cookie) error {
	request := struct {
		Cookie *types.Cookie `json:"cookie"`
	}{cookie}

	return c.Session.Execute("cookie", "POST", request)
}

func (c *Client) DeleteCookie(cookieName string) error {
	return c.Session.Execute("cookie/"+cookieName, "DELETE", nil)
}

func (c *Client) DeleteCookies() error {
	return c.Session.Execute("cookie", "DELETE", nil)
}

func (c *Client) GetScreenshot() ([]byte, error) {
	var base64Image string

	if err := c.Session.Execute("screenshot", "GET", nil, &base64Image); err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(base64Image)
}

func (c *Client) GetURL() (string, error) {
	var url string
	if err := c.Session.Execute("url", "GET", nil, &url); err != nil {
		return "", err
	}

	return url, nil
}

func (c *Client) SetURL(url string) error {
	request := struct {
		URL string `json:"url"`
	}{url}

	return c.Session.Execute("url", "POST", request)
}

func (c *Client) GetTitle() (string, error) {
	var title string
	if err := c.Session.Execute("title", "GET", nil, &title); err != nil {
		return "", err
	}

	return title, nil
}

func (c *Client) GetSource() (string, error) {
	var source string
	if err := c.Session.Execute("source", "GET", nil, &source); err != nil {
		return "", err
	}

	return source, nil
}

func (c *Client) DoubleClick() error {
	return c.Session.Execute("doubleclick", "POST", nil)
}

func (c *Client) MoveTo(element types.Element, point types.Point) error {
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

	return c.Session.Execute("moveto", "POST", request)
}

func (c *Client) Execute(body string, arguments []interface{}, result interface{}) error {
	request := struct {
		Script string        `json:"script"`
		Args   []interface{} `json:"args"`
	}{body, arguments}

	if err := c.Session.Execute("execute", "POST", request, result); err != nil {
		return err
	}

	return nil
}

func (c *Client) Forward() error {
	return c.Session.Execute("forward", "POST", nil)
}

func (c *Client) Back() error {
	return c.Session.Execute("back", "POST", nil)
}

func (c *Client) Refresh() error {
	return c.Session.Execute("refresh", "POST", nil)
}
