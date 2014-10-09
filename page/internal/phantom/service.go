package phantom

import (
	"encoding/json"
	"errors"
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
	Address string
	Timeout time.Duration
	process *os.Process
}

func (s *Service) Start() error {
	if _, err := exec.LookPath("phantomjs"); err != nil {
		return errors.New("phantomjs not found")
	}

	command := exec.Command("phantomjs", fmt.Sprintf("--webdriver=%s", s.Address))
	command.Start()
	s.process = command.Process

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
		return errors.New("phantomjs webdriver failed to start")
	case <-startedChan:
		return nil
	}
}

func (s *Service) Stop() {
	s.process.Signal(syscall.SIGINT)
	s.process.Wait()
	s.process = nil
}

func (s *Service) CreateSession() (*Session, error) {
	if s.process == nil {
		return nil, errors.New("phantomjs not running")
	}

	client := &http.Client{}
	postBody := strings.NewReader(`{"desiredCapabilities": {} }`)
	request, _ := http.NewRequest("POST", fmt.Sprintf("http://%s/session", s.Address), postBody)

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	var sessionResponse struct{ SessionID string }

	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &sessionResponse)

	if sessionResponse.SessionID == "" {
		return nil, errors.New("phantomjs webdriver failed to return a session ID")
	}

	sessionURL := fmt.Sprintf("http://%s/session/%s", s.Address, sessionResponse.SessionID)
	return &Session{sessionURL}, nil
}
