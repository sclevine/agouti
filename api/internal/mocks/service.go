package mocks

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
