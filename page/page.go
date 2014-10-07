package page

import "github.com/sclevine/agouti/webdriver"

type Page interface {
	Navigate(url string) Page
	SetCookie(cookie webdriver.Cookie) Page
	Selection
}

type page struct {
	driver driver
	fail   func(message string, callerSkip ...int)
}

type driver interface {
	Navigate(url string) error
	GetElements(selector string) ([]webdriver.Element, error)
	SetCookie(cookie *webdriver.Cookie) error
}

type callable interface {
	Call(selection Selection)
}

func NewPage(driver driver, fail func(message string, callerSkip ...int)) Page {
	return &page{driver, fail}
}

func (p *page) Navigate(url string) Page {
	if err := p.driver.Navigate(url); err != nil {
		p.fail(err.Error(), 1)
	}

	return p
}

func (p *page) SetCookie(cookie webdriver.Cookie) Page {
	if err := p.driver.SetCookie(&cookie); err != nil {
		p.fail(err.Error(), 1)
	}

	return p
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

func (p *page) Click() {
	p.body().Click()
}

func (p *page) body() *selection {
	return &selection{[]string{"body"}, p}
}
