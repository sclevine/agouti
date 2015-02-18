package mocks

import "github.com/sclevine/agouti/api"

type WebDriver struct {
	OpenCall struct {
		Desired       []api.Capabilities
		ReturnSession *api.Session
		Err           error
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

func (w *WebDriver) Open(desired ...api.Capabilities) (*api.Session, error) {
	return w.OpenCall.ReturnSession, w.OpenCall.Err
}

func (w *WebDriver) Start() error {
	w.StartCall.Called = true
	return w.StartCall.Err
}

func (w *WebDriver) Stop() error {
	w.StopCall.Called = true
	return w.StopCall.Err
}
