package phantom

import (
	"fmt"
	"os/exec"
	"os"
	"syscall"
	"errors"
	"net/http"
	"time"
	"strings"
	"io/ioutil"
	"encoding/json"
)

const PHANTOM_HOST = "127.0.0.1"
const PHANTOM_PORT = 8910

type Phantom struct {
	Timeout time.Duration
	process *os.Process
}

func (p *Phantom) Start() error {
	_, err := exec.LookPath("phantomjs")
	if err != nil {
		return errors.New("phantomjs not found")
	}

	command := exec.Command("phantomjs", fmt.Sprintf("--webdriver=%s:%d", PHANTOM_HOST, PHANTOM_PORT))
	command.Start()
	p.process = command.Process

	return p.waitForServer()
}

func (p *Phantom) waitForServer() error {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/status", PHANTOM_HOST, PHANTOM_PORT), nil)

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

func (p *Phantom) CreateSession() (sessionURL string, err error) {
	if p.process == nil {
		err = errors.New("phantomjs not running")
		return
	}

	client := &http.Client{}
	postBody := strings.NewReader(`{"desiredCapabilities": {} }`)
	request, _ := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/session", PHANTOM_HOST, PHANTOM_PORT), postBody)

    response, err := client.Do(request)

	if err != nil {
		return
	}

	var sessionResponse struct{SessionID string}

	body, _ := ioutil.ReadAll(response.Body)

	_ = json.Unmarshal(body, &sessionResponse)

	sessionURL = fmt.Sprintf("http://%s:%d/session/%s", PHANTOM_HOST, PHANTOM_PORT, sessionResponse.SessionID)
	return sessionURL, nil
}
