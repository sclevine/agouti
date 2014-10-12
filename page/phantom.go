package page

import (
	"fmt"
	"github.com/sclevine/agouti/page/internal/service"
	"github.com/sclevine/agouti/page/internal/webdriver"
	"net"
	"os/exec"
	"time"
)

var phantomService *service.Service

func StartPhantom() error {
	address, err := phantomFreeAddress()
	if err != nil {
		return fmt.Errorf("failed to locate a free port: %s", err)
	}

	desiredCapabilities := `{"desiredCapabilities": {} }`
	command := exec.Command("phantomjs", fmt.Sprintf("--webdriver=%s", address))
	phantomService = &service.Service{Address: address,
		Timeout:             3 * time.Second,
		ServiceName:         "phantomjs",
		Command:             command,
		DesiredCapabilities: desiredCapabilities}

	if err := phantomService.Start(); err != nil {
		return fmt.Errorf("failed to start PhantomJS: %s", err)
	}
	return nil
}

func phantomFreeAddress() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer listener.Close()
	return listener.Addr().String(), nil
}

func StopPhantom() {
	phantomService.Stop()
}

func PhantomPage() (*Page, error) {
	session, err := phantomService.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PhantomJS page: %s", err)
	}

	driver := &webdriver.Driver{session}
	return &Page{driver}, nil
}
