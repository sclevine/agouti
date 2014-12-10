package api

import (
	"encoding/base64"

	"errors"
	"github.com/sclevine/agouti/core/internal/api/element"
	"github.com/sclevine/agouti/core/internal/api/window"
	"github.com/sclevine/agouti/core/internal/types"
)

type Client struct {
	Session types.Session
}

func (c *Client) DeleteSession() error {
	return c.Session.Execute("", "DELETE", nil)
}

func (c *Client) GetElement(selector types.Selector) (types.Element, error) {
	var result struct{ Element string }

	if err := c.Session.Execute("element", "POST", selector, &result); err != nil {
		return nil, err
	}

	return &element.Element{ID: result.Element, Session: c.Session}, nil
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

func (c *Client) GetActiveElement() (types.Element, error) {
	var result struct{ Element string }

	if err := c.Session.Execute("element/active", "POST", nil, &result); err != nil {
		return nil, err
	}

	return &element.Element{ID: result.Element, Session: c.Session}, nil
}

func (c *Client) GetWindow() (types.Window, error) {
	var windowID string
	if err := c.Session.Execute("window_handle", "GET", nil, &windowID); err != nil {
		return nil, err
	}
	return &window.Window{ID: windowID, Session: c.Session}, nil
}

func (c *Client) SetCookie(cookie interface{}) error {
	request := struct {
		Cookie interface{} `json:"cookie"`
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

func (c *Client) MoveTo(region types.Element, point types.Point) error {
	request := map[string]interface{}{}

	if region != nil {
		// TODO: return error if not an element
		request["element"] = region.(*element.Element).ID
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

func (c *Client) Frame(frame types.Element) error {
	var elementID interface{}

	switch frame := frame.(type) {
	case nil:
		elementID = nil
	case *element.Element:
		elementID = struct {
			Element string `json:"ELEMENT"`
		}{frame.ID}
	default:
		return errors.New("frame must be an element")
	}

	request := struct {
		ID interface{} `json:"id"`
	}{elementID}

	return c.Session.Execute("frame", "POST", request)
}

func (c *Client) FrameParent() error {
	return c.Session.Execute("frame/parent", "POST", nil)
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

func (c *Client) GetAlertText() (string, error) {
	var text string
	if err := c.Session.Execute("alert_text", "GET", nil, &text); err != nil {
		return "", err
	}
	return text, nil
}

func (c *Client) SetAlertText(text string) error {
	request := struct {
		Text string `json:"text"`
	}{text}
	return c.Session.Execute("alert_text", "POST", request)
}

func (c *Client) NewLogs(logType string) ([]types.Log, error) {
	request := struct {
		Type string `json:"type"`
	}{logType}

	var logs []types.Log
	if err := c.Session.Execute("log", "POST", request, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

func (c *Client) GetLogTypes() ([]string, error) {
	var types []string
	if err := c.Session.Execute("log/types", "GET", nil, &types); err != nil {
		return nil, err
	}
	return types, nil
}
