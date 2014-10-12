package page

import (
	"fmt"
	"github.com/sclevine/agouti/page/internal/service"
	"github.com/sclevine/agouti/page/internal/webdriver"
	"time"
)

var phantomService *service.Service

func StartPhantom() error {
	address, err := freeAddress()
	if err != nil {
		return fmt.Errorf("failed to locate a free port: %s", err)
	}

	url := fmt.Sprintf("http://%s", address)
	command := []string{"phantomjs", fmt.Sprintf("--webdriver=%s", address)}
	phantomService = &service.Service{URL: url, Timeout: 3 * time.Second, Command: command}

	if err := phantomService.Start(); err != nil {
		return fmt.Errorf("failed to start PhantomJS: %s", err)
	}

	return nil
}

func StopPhantom() {
	phantomService.Stop()
}

func PhantomPage() (*Page, error) {
	capabilites := &service.Capabilities{}
	session, err := phantomService.CreateSession(capabilites)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PhantomJS page: %s", err)
	}

	driver := &webdriver.Driver{session}
	return &Page{driver}, nil
}
