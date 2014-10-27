package mocks

import (
	"github.com/sclevine/agouti/core/internal/session"
)

type Service struct {
	StartCall struct {
		Called bool
		Err    error
	}

	StopCall struct {
		Called bool
	}

	CreateSessionCall struct {
		Capabilities  session.JSONable
		ReturnSession *session.Session
		Err           error
	}
}

func (s *Service) Start() error {
	s.StartCall.Called = true
	return s.StartCall.Err
}

func (s *Service) Stop() {
	s.StopCall.Called = true
}

func (s *Service) CreateSession(capabilities session.JSONable) (*session.Session, error) {
	s.CreateSessionCall.Capabilities = capabilities
	return s.CreateSessionCall.ReturnSession, s.CreateSessionCall.Err
}
