package page

import "github.com/sclevine/agouti/webdriver"

type Page interface {
	Selection
}

type page struct {
	driver driver
	fail   func(message string, callerSkip ...int)
}

type driver interface {
	Navigate(url string) error
	GetElements(selector string) ([]webdriver.Element, error)
}

type callable interface {
	Call(selection Selection)
}

func NewPage(driver driver, fail func(message string, callerSkip ...int)) Page {
	return &page{driver, fail}
}

func (p *page) Within(selector string, bodies ...callable) Selection {
	firstSelection := &selection{[]string{selector}, p}
	for _, body := range bodies {
		body.Call(firstSelection)
	}
	return firstSelection
}

func (p *page) Selector() string {
	return p.body().Selector()
}

func (p *page) ShouldContainText(text string) {
	p.body().ShouldContainText(text)
}

func (p *page) body() *selection {
	return &selection{[]string{"body"}, p}
}
