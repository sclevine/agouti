package page

import (
	"github.com/sclevine/agouti/webdriver"
	"github.com/sclevine/agouti"
)

type Page struct {
	Driver driver
}

type driver interface {
	Navigate(url string) error
	GetElements(selector string) ([]*webdriver.Element, error)
}

func (p *Page) Within(selector string, bodies ...func(agouti.Selection)) agouti.Selection {
	firstSelection := &selection{[]string{selector}, p}
	for _, body := range bodies {
		body(firstSelection)
	}
	return firstSelection
}

func (p *Page) ShouldContainText(text string) {
	p.body().ShouldContainText(text)
}

func (p *Page) body() *selection {
	return &selection{[]string{"body"}, p}
}
