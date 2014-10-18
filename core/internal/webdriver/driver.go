package webdriver

import (
	"errors"
	"fmt"
	"github.com/sclevine/agouti/core/internal/api"
	"github.com/sclevine/agouti/core/internal/page"
	"github.com/sclevine/agouti/core/internal/session"
	"github.com/sclevine/agouti/core/internal/types"
)

type Driver struct {
	Service service
	pages   []types.Page
}

type service interface {
	Start() error
	Stop()
	CreateSession(capabilities map[string]interface{}) (*session.Session, error)
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

func (d *Driver) Page(browserName ...string) (types.Page, error) {
	capabilites := map[string]interface{}{}
	if len(browserName) == 1 {
		capabilites["browserName"] = browserName[0]
	} else if len(browserName) > 1 {
		return nil, errors.New("too many arguments")
	}

	pageSession, err := d.Service.CreateSession(capabilites)
	if err != nil {
		return nil, fmt.Errorf("failed to generate page: %s", err)
	}

	pageClient := &api.Client{Session: pageSession}
	newPage := &page.Page{Client: pageClient}
	d.pages = append(d.pages, newPage)
	return newPage, nil
}
