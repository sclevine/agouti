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

type Phantom struct {
	Host    string
	Port    int
	Timeout time.Duration
	process *os.Process
}

func (p *Phantom) Start() error {
	if _, err := exec.LookPath("phantomjs"); err != nil {
		return errors.New("phantomjs not found")
	}

	command := exec.Command("phantomjs", fmt.Sprintf("--webdriver=%s:%d", p.Host, p.Port))
	command.Start()
	p.process = command.Process

	return p.waitForServer()
}

func (p *Phantom) waitForServer() error {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/status", p.Host, p.Port), nil)

	timeoutChan := time.After(p.Timeout)
	failedChan := make(chan struct{}, 1)
	startedChan := make(chan struct{})

	go func() {
		_, err := client.Do(request)
		for err != nil {
			select {
			case <-failedChan:
				return
			default:
				_, err = client.Do(request)
			}
		}
		startedChan <- struct{}{}
	}()

	select {
	case <-timeoutChan:
		failedChan <- struct{}{}
		p.Stop()
		return errors.New("phantomjs webdriver failed to start")
	case <-startedChan:
		return nil
	}
}

func (p *Phantom) Stop() {
	p.process.Signal(syscall.SIGINT)
	p.process = nil
}

func (p *Phantom) CreateSession() (*Session, error) {
	if p.process == nil {
		return nil, errors.New("phantomjs not running")
	}

	client := &http.Client{}
	postBody := strings.NewReader(`{"desiredCapabilities": {} }`)
	request, _ := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/session", p.Host, p.Port), postBody)

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

	sessionURL := fmt.Sprintf("http://%s:%d/session/%s", p.Host, p.Port, sessionResponse.SessionID)
	return &Session{sessionURL}, nil
}
