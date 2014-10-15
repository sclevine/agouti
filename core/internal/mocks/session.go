package mocks

import "encoding/json"

type Session struct {
	ExecuteCall struct {
		Endpoint string
		Method   string
		BodyJSON []byte
		Result   string
		Err      error
	}

	DestroyCall struct {
		Err error
	}
}

func (s *Session) Execute(endpoint, method string, body, result interface{}) error {
	s.ExecuteCall.Endpoint = endpoint
	s.ExecuteCall.Method = method
	s.ExecuteCall.BodyJSON, _ = json.Marshal(body)
	json.Unmarshal([]byte(s.ExecuteCall.Result), result)
	return s.ExecuteCall.Err
}

func (s *Session) Destroy() error {
	return s.DestroyCall.Err
}
