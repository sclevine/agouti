package mocks

import "encoding/json"

type Session struct {
	Endpoint string
	Method   string
	BodyJSON []byte
	Result   string
	Err      error
}

func (s *Session) Execute(endpoint, method string, body, result interface{}) error {
	s.Endpoint = endpoint
	s.Method = method
	s.BodyJSON, _ = json.Marshal(body)
	json.Unmarshal([]byte(s.Result), result)
	return s.Err
}
