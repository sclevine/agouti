package core

import (
	"time"

	"github.com/sclevine/agouti/core/internal/service"
)

// Chrome returns an instance of a ChromeDriver WebDriver.
// DEPRECATED: Use ChromeDriver() instead.
func Chrome() (WebDriver, error) {
	chrome := &service.Service{
		URLTemplate: "http://{{.Address}}",
		CmdTemplate: []string{"chromedriver", "--silent", "--port={{.Port}}"},
		Timeout:     5 * time.Second,
	}
	return &driver{service: chrome}, nil
}
