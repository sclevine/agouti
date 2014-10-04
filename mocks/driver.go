package mocks

import "github.com/sclevine/agouti/webdriver"

type Driver struct {
	Navigate struct {
		URL string
		Err error
	}

	GetElements struct {
		Selector string
		ReturnElements []webdriver.Element
		Err error
	}

	SetCookie struct {
		Cookie *webdriver.Cookie
		Err error
	}
}

func (d* Driver) Navigate(url string) error {
	d.Navigate.URL = url
	return d.Navigate.Err
}

func (d *Driver) GetElements(selector string) ([]webdriver.Element, error) {
	d.GetElements.Selector = selector
	return d.GetElements.ReturnElements, d.GetElements.Err
}

func (d *Driver) SetCookie(cookie *webdriver.Cookie) error {
	d.SetCookie.Cookie = cookie
	return d.SetCookie.Err
}
