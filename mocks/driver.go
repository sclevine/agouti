package mocks

import (
	"github.com/sclevine/agouti/webdriver"
)

type Driver struct {
	NavigateCall struct {
		URL string
		Err error
	}

	GetElementsCall struct {
		Selector       string
		ReturnElements []webdriver.Element
		Err            error
	}

	GetWindowCall struct {
	    ReturnWindow webdriver.Window
		Err error
	}

	SetCookieCall struct {
		Cookie *webdriver.Cookie
		Err    error
	}

	GetURLCall struct {
		ReturnURL string
		Err       error
	}
}

func (d *Driver) Navigate(url string) error {
	d.NavigateCall.URL = url
	return d.NavigateCall.Err
}

func (d *Driver) GetElements(selector string) ([]webdriver.Element, error) {
	d.GetElementsCall.Selector = selector
	return d.GetElementsCall.ReturnElements, d.GetElementsCall.Err
}

func(d *Driver) GetWindow() (webdriver.Window, error) {
	return d.GetWindowCall.ReturnWindow, d.GetWindowCall.Err
}

func (d *Driver) SetCookie(cookie *webdriver.Cookie) error {
	d.SetCookieCall.Cookie = cookie
	return d.SetCookieCall.Err
}

func (d *Driver) GetURL() (string, error) {
	return d.GetURLCall.ReturnURL, d.GetURLCall.Err
}
