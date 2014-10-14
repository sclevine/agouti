package mocks

import "github.com/sclevine/agouti/core/internal/webdriver"

type Driver struct {
	GetElementsCall struct {
		Selector       string
		ReturnElements []webdriver.Element
		Err            error
	}

	GetWindowCall struct {
		ReturnWindow webdriver.Window
		Err          error
	}

	GetScreenshotCall struct {
		ReturnImage []byte
		Err         error
	}

	SetCookieCall struct {
		Cookie *webdriver.Cookie
		Err    error
	}

	DeleteCookieCall struct {
		Name string
		Err  error
	}

	DeleteCookiesCall struct {
		WasCalled bool
		Err       error
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
}

func (d *Driver) GetElements(selector string) ([]webdriver.Element, error) {
	d.GetElementsCall.Selector = selector
	return d.GetElementsCall.ReturnElements, d.GetElementsCall.Err
}

func (d *Driver) GetWindow() (webdriver.Window, error) {
	return d.GetWindowCall.ReturnWindow, d.GetWindowCall.Err
}

func (d *Driver) GetScreenshot() ([]byte, error) {
	return d.GetScreenshotCall.ReturnImage, d.GetScreenshotCall.Err
}

func (d *Driver) SetCookie(cookie *webdriver.Cookie) error {
	d.SetCookieCall.Cookie = cookie
	return d.SetCookieCall.Err
}

func (d *Driver) DeleteCookie(name string) error {
	d.DeleteCookieCall.Name = name
	return d.DeleteCookieCall.Err
}

func (d *Driver) DeleteCookies() error {
	d.DeleteCookiesCall.WasCalled = true
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
