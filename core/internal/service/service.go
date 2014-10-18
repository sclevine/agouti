package service

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/sclevine/agouti/core/internal/session"
)

type Service struct {
	URL     string
	Timeout time.Duration
	Command []string
	process *os.Process
}

func (s *Service) name() string {
	return s.Command[0]
}

func (s *Service) Start() error {
	if s.process != nil {
		return fmt.Errorf("%s is already running", s.name())
	}

	command := exec.Command(s.name(), s.Command[1:]...)

	if err := command.Start(); err != nil {
		return fmt.Errorf("unable to run %s: %s", s.name(), err)
	}

	s.process = command.Process

	return s.waitForServer()
}

func (s *Service) waitForServer() error {
	timeoutChan := time.After(s.Timeout)
	failedChan := make(chan struct{}, 1)
	startedChan := make(chan struct{})

	go func() {
		up := s.checkStatus()
		for !up {
			select {
			case <-failedChan:
				return
			default:
				time.Sleep(500 * time.Millisecond)
				up = s.checkStatus()
			}
		}
		startedChan <- struct{}{}
	}()

	select {
	case <-timeoutChan:
		failedChan <- struct{}{}
		s.Stop()
		return fmt.Errorf("%s failed to start", s.name())
	case <-startedChan:
		return nil
	}
}

func (s *Service) checkStatus() bool {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/status", s.URL), nil)
	response, err := client.Do(request)
	if err == nil && response.StatusCode == 200 {
		return true
	}
	return false
}

func (s *Service) Stop() {
	if s.process == nil {
		return
	}
	s.process.Signal(syscall.SIGINT)
	s.process.Wait()
	s.process = nil
}

func (s *Service) CreateSession(capabilities map[string]interface{}) (*session.Session, error) {
	if s.process == nil {
		return nil, fmt.Errorf("%s not running", s.name())
	}

	return session.Open(s.URL, capabilities)
}
