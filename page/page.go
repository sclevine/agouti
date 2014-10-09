package page

import (
	"github.com/sclevine/agouti/webdriver"
	"time"
)

type Page interface {
	Navigate(url string) Page
	SetCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) Page
	URL() string
	Size(height, width int)
	Selection
}

type page struct {
	driver driver
	failer failer
}

type failer interface {
	Fail(message string)
	Down()
	Up()
	Async()
	Sync()
	Reset()
}

type driver interface {
	Navigate(url string) error
	GetElements(selector string) ([]webdriver.Element, error)
	GetWindow() (webdriver.Window, error)
	SetCookie(cookie *webdriver.Cookie) error
	GetURL() (string, error)
}

type callable interface {
	Call(selection Selection)
}

func NewPage(driver driver, failer failer) Page {
	return &page{driver, failer}
}

func (p *page) Navigate(url string) Page {
	p.failer.Down()
	if err := p.driver.Navigate(url); err != nil {
		p.failer.Fail("Failed to navigate: " + err.Error())
	}
	p.failer.Up()
	return p
}

func (p *page) SetCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) Page {
	cookie := webdriver.Cookie{name, value, path, domain, secure, httpOnly, expiry}
	p.failer.Down()
	if err := p.driver.SetCookie(&cookie); err != nil {
		p.failer.Fail("Failed to set cookie: " + err.Error())
	}
	p.failer.Up()
	return p
}

func (p *page) URL() string {
	p.failer.Down()
	url, err := p.driver.GetURL()
	if err != nil {
		p.failer.Fail("Failed to retrieve URL: " + err.Error())
	}
	p.failer.Up()
	return url
}

func (p *page) Size(height, width int) {
	p.failer.Down()

	window, err := p.driver.GetWindow()
	if err != nil {
		p.failer.Fail("Failed to get a window: " + err.Error())
	}

	p.failer.Up()

	err = window.SetSize(640, 480)
	if err != nil {
		p.failer.Fail("Failed to re-size the window: " + err.Error())
	}
}

func (p *page) Should() FinalSelection {
	return p.body()
}

func (p *page) ShouldNot() FinalSelection {
	body := p.body()
	body.invert = true
	return body
}

func (p *page) ShouldEventually(timing ...time.Duration) FinalSelection {
	return p.body().ShouldEventually(timing...)
}

func (p *page) Within(selector string, bodies ...callable) Selection {
	firstSelection := &selection{p.driver, p.failer, []string{selector}, false}
	for _, body := range bodies {
		body.Call(firstSelection)
	}
	return firstSelection
}

func (p *page) Selector() string {
	return p.body().Selector()
}

func (p *page) Click() {
	p.failer.Down()
	p.body().Click()
	p.failer.Up()
}

func (p *page) body() *selection {
	return &selection{p.driver, p.failer, []string{"body"}, false}
}
