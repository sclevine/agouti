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
}

func (s *Session) Execute(endpoint, method string, body interface{}, result ...interface{}) error {
	s.ExecuteCall.Endpoint = endpoint
	s.ExecuteCall.Method = method
	s.ExecuteCall.BodyJSON, _ = json.Marshal(body)
	if len(result) > 0 {
		json.Unmarshal([]byte(s.ExecuteCall.Result), result[0])
	}
	return s.ExecuteCall.Err
}
