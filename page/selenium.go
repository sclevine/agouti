package page

import (
	"fmt"
	"github.com/sclevine/agouti/page/internal/service"
	"github.com/sclevine/agouti/page/internal/webdriver"
	"net"
	"os/exec"
	"strings"
	"time"
)

var seleniumService *service.Service

func StartSelenium(browserType string) error {
	addressString, err := seleniumFreeAddress()
	if err != nil {
		return fmt.Errorf("failed to locate a free port: %s", err)
	}

	port := strings.Split(addressString, ":")[1]
	address := fmt.Sprintf("127.0.0.1:%d/wd/hub", port)
	desiredCapabilities := fmt.Sprintf(`{"desiredCapabilities": {"browserName: "%s"}}`, browserType)
	command := exec.Command("selenium-server", fmt.Sprintf("-port %d", port))

	seleniumService = &service.Service{Address: address,
		Timeout:             5 * time.Second,
		ServiceName:         "selenium-server",
		Command:             command,
		DesiredCapabilities: desiredCapabilities}

	if err := seleniumService.Start(); err != nil {
		return fmt.Errorf("failed to start Selenium: %s", err)
	}
	return nil
}

func seleniumFreeAddress() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer listener.Close()
	return listener.Addr().String(), nil
}

func StopSelenium() {
	seleniumService.Stop()
}

func SeleniumPage() (*Page, error) {
	session, err := seleniumService.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Selenium page: %s", err)
	}

	driver := &webdriver.Driver{session}
	return &Page{driver}, nil
}
