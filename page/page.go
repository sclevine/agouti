package page

import "github.com/sclevine/agouti/webdriver"

type Page struct {
	Driver driver
	Fail   func(message string, callerSkip ...int)

}

type driver interface {
	Navigate(url string) error
	GetElements(selector string) ([]webdriver.Element, error)
}

type callable interface {
	Call(selection Selection)
}

func (p *Page) Within(selector string, bodies ...callable) Selection {
	firstSelection := &selection{[]string{selector}, p}
	for _, body := range bodies {
		body.Call(firstSelection)
	}
	return firstSelection
}

func (p *Page) Selector() string {
	return p.body().Selector()
}

func (p *Page) ShouldContainText(text string) {
	p.body().ShouldContainText(text)
}

func (p *Page) body() *selection {
	return &selection{[]string{"body"}, p}
}
