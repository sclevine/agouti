package page

import (
	"fmt"
	"github.com/sclevine/agouti/page/internal/service"
	"github.com/sclevine/agouti/page/internal/session"
	"github.com/sclevine/agouti/page/internal/webdriver"
	"strings"
	"time"
)

var (
	chromeService *service.Service
	chromeSessions []*session.Session
)

func StartChrome() error {
	address, err := freeAddress()
	if err != nil {
		return fmt.Errorf("failed to locate a free port: %s", err)
	}

	port := strings.SplitN(address, ":", 2)[1]
	url := fmt.Sprintf("http://%s", address)
	command := []string{"chromedriver", "--silent", "--port=" + port}
	chromeService = &service.Service{URL: url, Timeout: 5 * time.Second, Command: command}

	if err := chromeService.Start(); err != nil {
		return fmt.Errorf("failed to start Chrome: %s", err)
	}

	return nil
}

func StopChrome() {
	for _, session := range chromeSessions {
		session.Destroy()
	}
	chromeService.Stop()
}

func ChromePage() (*Page, error) {
	capabilites := &service.Capabilities{}
	session, err := chromeService.CreateSession(capabilites)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Chrome page: %s", err)
	}

	chromeSessions = append(chromeSessions, session)
	driver := &webdriver.Driver{session}
	return &Page{driver}, nil
}
