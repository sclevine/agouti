package page

import (
	"fmt"
	"github.com/sclevine/agouti/page/internal/phantom"
	"github.com/sclevine/agouti/page/internal/webdriver"
	"net"
	"time"
)

var phantomService *phantom.Service

func StartPhantom() error {
	address, err := freeAddress()
	if err != nil {
		return fmt.Errorf("failed to locate a free port: %s", err)
	}

	phantomService = &phantom.Service{Address: address, Timeout: 3 * time.Second}
	if err := phantomService.Start(); err != nil {
		return fmt.Errorf("failed to start PhantomJS: %s", err)
	}
	return nil
}

func freeAddress() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer listener.Close()
	return listener.Addr().String(), nil
}

func StopPhantom(startErr ...error) {
	if len(startErr) == 0 || startErr[0] == nil {
		phantomService.Stop()
	}
}

func PhantomPage() (*Page, error) {
	session, err := phantomService.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PhantomJS page: ", err.Error())
	}

	driver := &webdriver.Driver{session}
	return &Page{driver}, nil
}
