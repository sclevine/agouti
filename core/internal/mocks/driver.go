package mocks

import (
	"encoding/json"
	"github.com/sclevine/agouti/core/internal/webdriver/types"
)

type Driver struct {
	GetElementsCall struct {
		Selector       types.Selector
		ReturnElements []types.Element
		Err            error
	}

	GetWindowCall struct {
		ReturnWindow types.Window
		Err          error
	}

	GetScreenshotCall struct {
		ReturnImage []byte
		Err         error
	}

	SetCookieCall struct {
		Cookie *types.Cookie
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

	DoubleClickCall struct {
		Called bool
		Err    error
	}

	MoveToCall struct {
		Element types.Element
		Point   types.Point
		Err     error
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
}

func (d *Driver) GetElements(selector types.Selector) ([]types.Element, error) {
	d.GetElementsCall.Selector = selector
	return d.GetElementsCall.ReturnElements, d.GetElementsCall.Err
}

func (d *Driver) GetWindow() (types.Window, error) {
	return d.GetWindowCall.ReturnWindow, d.GetWindowCall.Err
}

func (d *Driver) GetScreenshot() ([]byte, error) {
	return d.GetScreenshotCall.ReturnImage, d.GetScreenshotCall.Err
}

func (d *Driver) SetCookie(cookie *types.Cookie) error {
	d.SetCookieCall.Cookie = cookie
	return d.SetCookieCall.Err
}

func (d *Driver) DeleteCookie(name string) error {
	d.DeleteCookieCall.Name = name
	return d.DeleteCookieCall.Err
}

func (d *Driver) DeleteCookies() error {
	d.DeleteCookiesCall.Called = true
	return d.DeleteCookiesCall.Err
}

func (d *Driver) GetURL() (string, error) {
	return d.GetURLCall.ReturnURL, d.GetURLCall.Err
}

func (d *Driver) SetURL(url string) error {
	d.SetURLCall.URL = url
	return d.SetURLCall.Err
}

func (d *Driver) GetTitle() (string, error) {
	return d.GetTitleCall.ReturnTitle, d.GetTitleCall.Err
}

func (d *Driver) DoubleClick() error {
	d.DoubleClickCall.Called = true
	return d.DoubleClickCall.Err
}

func (d *Driver) MoveTo(element types.Element, point types.Point) error {
	d.MoveToCall.Element = element
	d.MoveToCall.Point = point
	return d.MoveToCall.Err
}

func (d *Driver) Execute(body string, arguments []interface{}, result interface{}) error {
	d.ExecuteCall.Body = body
	d.ExecuteCall.Arguments = arguments
	json.Unmarshal([]byte(d.ExecuteCall.Result), result)
	return d.ExecuteCall.Err
}

func (d *Driver) Forward() error {
	d.ForwardCall.Called = true
	return d.ForwardCall.Err
}

func (d *Driver) Back() error {
	d.BackCall.Called = true
	return d.BackCall.Err
}

func (d *Driver) Refresh() error {
	d.RefreshCall.Called = true
	return d.RefreshCall.Err
}
