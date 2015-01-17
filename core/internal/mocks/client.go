package mocks

import (
	"encoding/json"

	"github.com/sclevine/agouti/core/internal/api"
)

type Client struct {
	GetElementCall struct {
		Selector      api.Selector
		ReturnElement *api.Element
		Err           error
	}

	GetElementsCall struct {
		Selector       api.Selector
		ReturnElements []*api.Element
		Err            error
	}

	GetActiveElementCall struct {
		ReturnElement *api.Element
		Err           error
	}

	DeleteSessionCall struct {
		Called bool
		Err    error
	}

	GetWindowCall struct {
		ReturnWindow *api.Window
		Err          error
	}

	GetWindowsCall struct {
		ReturnWindows []*api.Window
		Err           error
	}

	SetWindowCall struct {
		Window *api.Window
		Err    error
	}

	SetWindowByNameCall struct {
		Name string
		Err  error
	}

	DeleteWindowCall struct {
		Called bool
		Err    error
	}

	GetScreenshotCall struct {
		ReturnImage []byte
		Err         error
	}

	SetCookieCall struct {
		Cookie interface{}
		Err    error
	}

	DeleteCookieCall struct {
		Name string
		Err  error
	}

	DeleteCookiesCall struct {
		Called bool
		Err    error
	}

	GetURLCall struct {
		ReturnURL string
		Err       error
	}

	SetURLCall struct {
		URL string
		Err error
	}

	GetTitleCall struct {
		ReturnTitle string
		Err         error
	}

	GetSourceCall struct {
		ReturnSource string
		Err          error
	}

	DoubleClickCall struct {
		Called bool
		Err    error
	}

	MoveToCall struct {
		Element *api.Element
		Point   api.Point
		Err     error
	}

	FrameCall struct {
		Frame *api.Element
		Err   error
	}

	FrameParentCall struct {
		Called bool
		Err    error
	}

	ExecuteCall struct {
		Body      string
		Arguments []interface{}
		Result    string
		Err       error
	}

	ForwardCall struct {
		Called bool
		Err    error
	}

	BackCall struct {
		Called bool
		Err    error
	}

	RefreshCall struct {
		Called bool
		Err    error
	}

	GetAlertTextCall struct {
		ReturnText string
		Err        error
	}

	SetAlertTextCall struct {
		Text string
		Err  error
	}

	NewLogsCall struct {
		LogType    string
		ReturnLogs []api.Log
		Err        error
	}

	GetLogTypesCall struct {
		ReturnTypes []string
		Err         error
	}
}

func (c *Client) DeleteSession() error {
	c.DeleteSessionCall.Called = true
	return c.DeleteSessionCall.Err
}

func (c *Client) GetElement(selector api.Selector) (*api.Element, error) {
	c.GetElementCall.Selector = selector
	return c.GetElementCall.ReturnElement, c.GetElementCall.Err
}

func (c *Client) GetElements(selector api.Selector) ([]*api.Element, error) {
	c.GetElementsCall.Selector = selector
	return c.GetElementsCall.ReturnElements, c.GetElementsCall.Err
}

func (c *Client) GetActiveElement() (*api.Element, error) {
	return c.GetActiveElementCall.ReturnElement, c.GetActiveElementCall.Err
}

func (c *Client) GetWindow() (*api.Window, error) {
	return c.GetWindowCall.ReturnWindow, c.GetWindowCall.Err
}

func (c *Client) GetWindows() ([]*api.Window, error) {
	return c.GetWindowsCall.ReturnWindows, c.GetWindowsCall.Err
}

func (c *Client) SetWindow(window *api.Window) error {
	c.SetWindowCall.Window = window
	return c.SetWindowCall.Err
}

func (c *Client) SetWindowByName(name string) error {
	c.SetWindowByNameCall.Name = name
	return c.SetWindowByNameCall.Err
}

func (c *Client) DeleteWindow() error {
	c.DeleteWindowCall.Called = true
	return c.DeleteWindowCall.Err
}

func (c *Client) GetScreenshot() ([]byte, error) {
	return c.GetScreenshotCall.ReturnImage, c.GetScreenshotCall.Err
}

func (c *Client) SetCookie(cookie interface{}) error {
	c.SetCookieCall.Cookie = cookie
	return c.SetCookieCall.Err
}

func (c *Client) DeleteCookie(name string) error {
	c.DeleteCookieCall.Name = name
	return c.DeleteCookieCall.Err
}

func (c *Client) DeleteCookies() error {
	c.DeleteCookiesCall.Called = true
	return c.DeleteCookiesCall.Err
}

func (c *Client) GetURL() (string, error) {
	return c.GetURLCall.ReturnURL, c.GetURLCall.Err
}

func (c *Client) SetURL(url string) error {
	c.SetURLCall.URL = url
	return c.SetURLCall.Err
}

func (c *Client) GetTitle() (string, error) {
	return c.GetTitleCall.ReturnTitle, c.GetTitleCall.Err
}

func (c *Client) GetSource() (string, error) {
	return c.GetSourceCall.ReturnSource, c.GetSourceCall.Err
}

func (c *Client) DoubleClick() error {
	c.DoubleClickCall.Called = true
	return c.DoubleClickCall.Err
}

func (c *Client) MoveTo(element *api.Element, point api.Point) error {
	c.MoveToCall.Element = element
	c.MoveToCall.Point = point
	return c.MoveToCall.Err
}

func (c *Client) Frame(frame *api.Element) error {
	c.FrameCall.Frame = frame
	return c.FrameCall.Err
}

func (c *Client) FrameParent() error {
	c.FrameParentCall.Called = true
	return c.FrameParentCall.Err
}

func (c *Client) Execute(body string, arguments []interface{}, result interface{}) error {
	c.ExecuteCall.Body = body
	c.ExecuteCall.Arguments = arguments
	json.Unmarshal([]byte(c.ExecuteCall.Result), result)
	return c.ExecuteCall.Err
}

func (c *Client) Forward() error {
	c.ForwardCall.Called = true
	return c.ForwardCall.Err
}

func (c *Client) Back() error {
	c.BackCall.Called = true
	return c.BackCall.Err
}

func (c *Client) Refresh() error {
	c.RefreshCall.Called = true
	return c.RefreshCall.Err
}

func (c *Client) GetAlertText() (string, error) {
	return c.GetAlertTextCall.ReturnText, c.GetAlertTextCall.Err
}

func (c *Client) SetAlertText(text string) error {
	c.SetAlertTextCall.Text = text
	return c.SetAlertTextCall.Err
}

func (c *Client) NewLogs(logType string) ([]api.Log, error) {
	c.NewLogsCall.LogType = logType
	return c.NewLogsCall.ReturnLogs, c.NewLogsCall.Err
}

func (c *Client) GetLogTypes() ([]string, error) {
	return c.GetLogTypesCall.ReturnTypes, c.GetLogTypesCall.Err
}
