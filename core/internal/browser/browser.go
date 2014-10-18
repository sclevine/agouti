package browser

import (
	"errors"
	"fmt"
	"github.com/sclevine/agouti/core/internal/page"
	"github.com/sclevine/agouti/core/internal/service"
	"github.com/sclevine/agouti/core/internal/session"
	"github.com/sclevine/agouti/core/internal/types"
	"github.com/sclevine/agouti/core/internal/webdriver"
)

type Browser struct {
	Service browserService
	pages   []types.Page
}

type browserService interface {
	Start() error
	Stop()
	CreateSession(capabilities *service.Capabilities) (*session.Session, error)
}

func (b *Browser) Start() error {
	if err := b.Service.Start(); err != nil {
		return fmt.Errorf("failed to start service: %s", err)
	}

	return nil
}

func (b *Browser) Stop() {
	for _, browserPage := range b.pages {
		browserPage.Destroy()
	}

	b.Service.Stop()
	return
}

func (b *Browser) Page(browserName ...string) (types.Page, error) {
	capabilites := &service.Capabilities{}
	if len(browserName) == 1 {
		capabilites.BrowserName = browserName[0]
	} else if len(browserName) > 1 {
		return nil, errors.New("too many arguments")
	}

	pageSession, err := b.Service.CreateSession(capabilites)
	if err != nil {
		return nil, fmt.Errorf("failed to generate page: %s", err)
	}

	pageDriver := &webdriver.Driver{Session: pageSession}
	newPage := &page.Page{Driver: pageDriver}
	b.pages = append(b.pages, newPage)
	return newPage, nil
}
