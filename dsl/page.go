package dsl

import (
	"github.com/onsi/ginkgo"
	"github.com/sclevine/agouti/page"
)

type Page interface {
	Navigate(url string) error
	SetCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) error
	DeleteCookie(name string) error
	ClearCookies() error
	URL() (string, error)
	Size(height, width int) error
	Screenshot(filename string) error
	page.Selection
}

func CreatePage() Page {
	page, err := page.PhantomPage()
	if err != nil {
		ginkgo.Fail(err.Error(), 1)
	}

	return page
}
