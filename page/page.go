package page

import (
	"github.com/sclevine/agouti/webdriver"
)

type Page interface {
	Navigate(url string) Page
	SetCookie(cookie webdriver.Cookie) Page
	URL() string
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
	GetURL() (string, error)
}

type callable interface {
	Call(selection Selection)
}

func NewPage(driver driver, fail func(message string, callerSkip ...int)) Page {
	return &page{driver, fail}
}

func (p *page) Navigate(url string) Page {
	if err := p.driver.Navigate(url); err != nil {
		p.fail("Failed to navigate: "+err.Error(), 1)
	}

	return p
}

func (p *page) SetCookie(cookie webdriver.Cookie) Page {
	if err := p.driver.SetCookie(&cookie); err != nil {
		p.fail("Failed to set cookie: "+err.Error(), 1)
	}

	return p
}

func (p *page) URL() string {
	url, err := p.driver.GetURL()
	if err != nil {
		p.fail("Failed to retrieve URL: "+err.Error(), 1)
	}
	return url
}

func (p *page) Should() FinalSelection {
	return p.body()
}

func (p *page) ShouldNot() FinalSelection {
	body := p.body()
	body.invert = true
	return body
}

func (p *page) ShouldEventually() FinalSelection { // TODO: test
	return &async{p.body()}
}

func (p *page) Within(selector string, bodies ...callable) Selection {
	firstSelection := &selection{[]string{selector}, p, false}
	for _, body := range bodies {
		body.Call(firstSelection)
	}
	return firstSelection
}

func (p *page) Selector() string {
	return p.body().Selector()
}

func (p *page) ContainText(text string) {
	p.body().ContainText(text)
}

func (p *page) Click() {
	p.body().Click()
}

func (p *page) body() *selection {
	return &selection{[]string{"body"}, p, false}
}
