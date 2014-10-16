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
	Service  browserService
	sessions []destroyable
}

type browserService interface {
	Start() error
	Stop()
	CreateSession(capabilities *service.Capabilities) (*session.Session, error)
}

type destroyable interface {
	Destroy() error
}

func (b *Browser) Start() error {
	if err := b.Service.Start(); err != nil {
		return fmt.Errorf("failed to start service: %s", err)
	}

	return nil
}

func (b *Browser) Stop() (nonFatal error) {
	for _, pageSession := range b.sessions {
		if err := pageSession.Destroy(); err != nil {
			nonFatal = errors.New("failed to destroy all running sessions")
		}
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

	b.sessions = append(b.sessions, pageSession)
	pageDriver := &webdriver.Driver{pageSession}
	return &page.Page{pageDriver}, nil
}
