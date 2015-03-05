package appium

import (
	"fmt"

	"github.com/sclevine/agouti"
)

type Device struct {
	*agouti.Page
	session deviceSession
}

type deviceSession interface {
	SetEndpoint(thing string) error
}

func (d *Device) DeviceMethod(thing string) error {
	if err := d.session.SetEndpoint(thing); err != nil {
		return fmt.Errorf("failed to do stuff: %s", err)
	}
	return nil
}

// override finder methods
func (d *Device) Find(selector string) *Selection {
	return &Selection{d.Page.Find(selector), d.session}
}
