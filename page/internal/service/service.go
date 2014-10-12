package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type Service struct {
	Address             string
	Timeout             time.Duration
	ServiceName         string
	Command             *exec.Cmd
	DesiredCapabilities string
	process             *os.Process
	sessionID           string
}

func (s *Service) Start() error {
	if s.process != nil {
		return fmt.Errorf("%s is already running", s.ServiceName)
	}

	if _, err := exec.LookPath(s.ServiceName); err != nil {
		return fmt.Errorf("%s binary not found", s.ServiceName)
	}

	s.Command.Start()
	s.process = s.Command.Process

	return s.waitForServer()
}

func (s *Service) waitForServer() error {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", fmt.Sprintf("http://%s/status", s.Address), nil)

	timeoutChan := time.After(s.Timeout)
	failedChan := make(chan struct{}, 1)
	startedChan := make(chan struct{})

	go func() {
		_, err := client.Do(request)
		for err != nil {
			select {
			case <-failedChan:
				return
			default:
				time.Sleep(500 * time.Millisecond)
				_, err = client.Do(request)
			}
		}
		startedChan <- struct{}{}
	}()

	select {
	case <-timeoutChan:
		failedChan <- struct{}{}
		s.Stop()
		return fmt.Errorf("%s webdriver failed to start", s.ServiceName)
	case <-startedChan:
		return nil
	}
}

func (s *Service) Stop() {
	if s.ServiceName == "selenium-server" {
		client := &http.Client{}
		request, _ := http.NewRequest("DELETE", fmt.Sprintf("http://%s/session/%s", s.Address, s.sessionID), nil)
		_, _ = client.Do(request)
	}
	s.process.Signal(syscall.SIGINT)
	s.process.Wait()
	s.process = nil
}

func (s *Service) CreateSession() (*Session, error) {
	if s.process == nil {
		return nil, fmt.Errorf("%s not running", s.ServiceName)
	}

	client := &http.Client{}
	postBody := strings.NewReader(s.DesiredCapabilities)
	request, _ := http.NewRequest("POST", fmt.Sprintf("http://%s/session", s.Address), postBody)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	var sessionResponse struct{ SessionID string }

	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &sessionResponse)

	if sessionResponse.SessionID == "" {
		return nil, fmt.Errorf("%s webdriver failed to return a session ID", s.ServiceName)
	}

	s.sessionID = sessionResponse.SessionID

	sessionURL := fmt.Sprintf("http://%s/session/%s", s.Address, sessionResponse.SessionID)
	return &Session{sessionURL}, nil
}
