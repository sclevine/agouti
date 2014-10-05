package page

import (
	"github.com/sclevine/agouti/webdriver"
)

type Page struct {
	Driver driver
}

type driver interface {
	Navigate(url string) error
	GetElements(selector string) ([]*webdriver.Element, error)
}


type SelectionFunc interface {
	Call(selection *PageSelection)
}

func (p *Page) Within(selector string, bodies ...SelectionFunc) *PageSelection {
	firstSelection := &PageSelection{[]string{selector}, p}
	for _, body := range bodies {
		body.Call(firstSelection)
	}
	return firstSelection
}

func (p *Page) ShouldContainText(text string) {
	p.body().ShouldContainText(text)
}

func (p *Page) body() *PageSelection {
	return &PageSelection{[]string{"body"}, p}
}
