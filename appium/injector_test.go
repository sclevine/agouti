package appium

import "github.com/bradbev/agouti"

func NewTestDevice(session mobileSession) *Device {
	return &Device{
		Page:    &agouti.Page{},
		session: session,
	}
}
