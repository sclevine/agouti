package service

import (
	"encoding/json"
	"fmt"
	"github.com/sclevine/agouti/core/internal/session"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type Service struct {
	URL     string
	Timeout time.Duration
	Command []string
	process *os.Process
}

type Capabilities struct {
	BrowserName string `json:"browserName,omitempty"`
	Version     string `json:"version,omitempty"`
	Platform    string `json:"platform,omitempty"`
}

func (s *Service) name() string {
	return s.Command[0]
}

func (s *Service) Start() error {
	if s.process != nil {
		return fmt.Errorf("%s is already running", s.name())
	}

	if _, err := exec.LookPath(s.name()); err != nil {
		return fmt.Errorf("%s binary not found", s.name())
	}

	command := exec.Command(s.name(), s.Command[1:]...)

	command.Start()
	s.process = command.Process

	return s.waitForServer()
}

func (s *Service) waitForServer() error {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/status", s.URL), nil)

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
		return fmt.Errorf("%s webdriver failed to start", s.name())
	case <-startedChan:
		return nil
	}
}

func (s *Service) Stop() {
	s.process.Signal(syscall.SIGINT)
	s.process.Wait()
	s.process = nil
}

func (s *Service) CreateSession(capabilities *Capabilities) (*session.Session, error) {
	if s.process == nil {
		return nil, fmt.Errorf("%s not running", s.name())
	}

	capabilitiesJSON, _ := json.Marshal(capabilities)
	desiredCapabilities := fmt.Sprintf(`{"desiredCapabilities": %s}`, capabilitiesJSON)
	postBody := strings.NewReader(desiredCapabilities)

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/session", s.URL), postBody)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	var sessionResponse struct{ SessionID string }

	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &sessionResponse)

	if sessionResponse.SessionID == "" {
		return nil, fmt.Errorf("%s webdriver failed to return a session ID", s.name())
	}

	sessionURL := fmt.Sprintf("%s/session/%s", s.URL, sessionResponse.SessionID)
	return &session.Session{URL: sessionURL}, nil
}
