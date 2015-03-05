package mocks

import "time"

type Service struct {
	URLCall struct {
		ReturnURL string
		Err       error
	}

	StartCall struct {
		Called bool
		Err    error
	}

	StopCall struct {
		Called bool
		Err    error
	}

	WaitForBootCall struct {
		Timeout time.Duration
		Err     error
	}
}

func (s *Service) URL() (string, error) {
	return s.URLCall.ReturnURL, s.URLCall.Err
}

func (s *Service) Start() error {
	s.StartCall.Called = true
	return s.StartCall.Err
}

func (s *Service) Stop() error {
	s.StopCall.Called = true
	return s.StopCall.Err
}

func (s *Service) WaitForBoot(timeout time.Duration) error {
	s.WaitForBootCall.Timeout = timeout
	return s.WaitForBootCall.Err
}

func (s *Service) Debug(state bool) {
}
