package agouti

import (
	"github.com/onsi/ginkgo"
	"github.com/sclevine/agouti/webdriver"
)

type Page struct {
	driver *webdriver.Driver
}

func Navigate(url string, cookies []string) *Page {
	session, err := phantomService.CreateSession()
	// TODO: test error
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	driver := &webdriver.Driver{session}

	driver.SetCookies(cookies)
	driver.Navigate(url)
	return &Page{driver}
}

func (p *Page) Within(selector string, bodies ...func(*Selection)) *Selection {
	selection := &Selection{[]string{selector}, p}
	for _, body := range bodies {
		body(selection)
	}
	return selection
}

func (p *Page) ShouldContainText(text string) {
	p.pageSelection().ShouldContainText(text)
}

func (p *Page) pageSelection() *Selection {
	return &Selection{[]string{"body"}, p}
}
