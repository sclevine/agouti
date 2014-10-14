package core

import (
	"fmt"
	"github.com/sclevine/agouti/core/internal/browser"
	"github.com/sclevine/agouti/core/internal/page"
	"github.com/sclevine/agouti/core/internal/selection"
	"github.com/sclevine/agouti/core/internal/service"
	"net"
	"strings"
	"time"
)

type Selection selection.Selection
type Page page.Page

type Browser interface {
	Start() error
	Stop() (nonFatal error)
	Page(browserName ...string) (page.Page, error)
}

func Chrome() (Browser, error) {
	address, err := freeAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to locate a free port: %s", err)
	}

	port := strings.SplitN(address, ":", 2)[1]
	url := fmt.Sprintf("http://%s", address)
	command := []string{"chromedriver", "--silent", "--port=" + port}
	service := &service.Service{URL: url, Timeout: 5 * time.Second, Command: command}

	return &browser.Browser{Service: service}, nil
}

func PhantomJS() (Browser, error) {
	address, err := freeAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to locate a free port: %s", err)
	}

	url := fmt.Sprintf("http://%s", address)
	command := []string{"phantomjs", fmt.Sprintf("--webdriver=%s", address)}
	service := &service.Service{URL: url, Timeout: 3 * time.Second, Command: command}

	return &browser.Browser{Service: service}, nil
}

func Selenium() (Browser, error) {
	address, err := freeAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to locate a free port: %s", err)
	}

	port := strings.SplitN(address, ":", 2)[1]
	url := fmt.Sprintf("http://%s/wd/hub", address)
	command := []string{"selenium-server", "-port", port}
	service := &service.Service{URL: url, Timeout: 5 * time.Second, Command: command}

	return &browser.Browser{Service: service}, nil
}

func freeAddress() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer listener.Close()
	return listener.Addr().String(), nil
}
