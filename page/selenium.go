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
	seleniumService  *service.Service
	seleniumSessions []*session.Session
)

func StartSelenium() error {
	address, err := freeAddress()
	if err != nil {
		return fmt.Errorf("failed to locate a free port: %s", err)
	}

	port := strings.SplitN(address, ":", 2)[1]
	url := fmt.Sprintf("http://%s/wd/hub", address)
	command := []string{"selenium-server", "-port", port}
	seleniumService = &service.Service{URL: url, Timeout: 5 * time.Second, Command: command}

	if err := seleniumService.Start(); err != nil {
		return fmt.Errorf("failed to start Selenium: %s", err)
	}

	return nil
}

func StopSelenium() {
	for _, session := range seleniumSessions {
		session.Destroy()
	}
	seleniumService.Stop()
}

func SeleniumPage(pageType string) (*Page, error) {
	capabilites := &service.Capabilities{BrowserName: pageType}
	session, err := seleniumService.CreateSession(capabilites)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Selenium page: %s", err)
	}

	seleniumSessions = append(seleniumSessions, session)
	driver := &webdriver.Driver{session}
	return &Page{driver}, nil
}
