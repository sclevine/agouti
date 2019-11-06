package service

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

/** A running service is a service that is already running
  * and does not need to be started.
  * It has a fixed Url and port.
 */
type RunningService struct {
	sync.RWMutex
	Url         string
}

func (s *RunningService) URL() string {
	return s.Url
}

func (s *RunningService) Start(debug bool) error {
	// Nothing to do, it is already started
	return nil
}

func (s *RunningService) Stop() error {
	s.Lock()
	s.Url = ""
	s.Unlock()

	return nil
}

func (s *RunningService) WaitForBoot(timeout time.Duration) error {
	timeoutChan := time.After(timeout)
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
		return errors.New("failed to start before timeout")
	case <-startedChan:
		return nil
	}
}

func (s *RunningService) checkStatus() bool {
	client := &http.Client{}
	s.Lock()
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/status", s.Url), nil)
	s.Unlock()
	response, err := client.Do(request)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		return true
	}
	return false
}
