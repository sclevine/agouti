package mocks

import "github.com/sclevine/agouti/core/internal/types"

type Service struct {
	StartCall struct {
		Called bool
		Err    error
	}

	StopCall struct {
		Called bool
	}

	CreateSessionCall struct {
		Capabilities  types.JSON
		ReturnSession types.Session
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

func (s *Service) CreateSession(capabilities types.JSON) (types.Session, error) {
	s.CreateSessionCall.Capabilities = capabilities
	return s.CreateSessionCall.ReturnSession, s.CreateSessionCall.Err
}
