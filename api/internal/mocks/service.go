package mocks

import "time"

type Service struct {
	URLCall struct {
		ReturnURL string
		Err       error
	}

	StartCall struct {
		Timeout time.Duration
		Err     error
	}

	StopCall struct {
		Called bool
		Err    error
	}
}

func (s *Service) URL() (string, error) {
	return s.URLCall.ReturnURL, s.URLCall.Err
}

func (s *Service) Start(timeout time.Duration) error {
	s.StartCall.Timeout = timeout
	return s.StartCall.Err
}

func (s *Service) Stop() error {
	s.StopCall.Called = true
	return s.StopCall.Err
}
