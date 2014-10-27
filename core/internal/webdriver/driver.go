package webdriver

import (
	"errors"
	"fmt"

	"github.com/sclevine/agouti/core/internal/api"
	"github.com/sclevine/agouti/core/internal/page"
	"github.com/sclevine/agouti/core/internal/session"
)

type Driver struct {
	Service Service
	pages   []*page.Page
}

type Service interface {
	Start() error
	Stop()
	CreateSession(capabilities session.JSONable) (*session.Session, error)
}

func (d *Driver) Page(capabilities ...session.JSONable) (*page.Page, error) {
	if len(capabilities) == 0 {
		capabilities = append(capabilities, session.Capabilities{})
	} else if len(capabilities) > 1 {
		return nil, errors.New("too many arguments")
	}

	pageSession, err := d.Service.CreateSession(capabilities[0])
	if err != nil {
		return nil, fmt.Errorf("failed to generate page: %s", err)
	}

	pageClient := &api.Client{Session: pageSession}
	newPage := &page.Page{Client: pageClient}
	d.pages = append(d.pages, newPage)
	return newPage, nil
}

func (d *Driver) Start() error {
	if err := d.Service.Start(); err != nil {
		return fmt.Errorf("failed to start service: %s", err)
	}

	return nil
}

func (d *Driver) Stop() {
	for _, openPage := range d.pages {
		openPage.Destroy()
	}

	d.Service.Stop()
	return
}
