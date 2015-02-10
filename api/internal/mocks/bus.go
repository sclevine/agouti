package mocks

import "encoding/json"

type Bus struct {
	SendCall struct {
		Endpoint string
		Method   string
		BodyJSON []byte
		Result   string
		Err      error
	}
}

func (b *Bus) Send(endpoint, method string, body interface{}, result ...interface{}) error {
	b.SendCall.Endpoint = endpoint
	b.SendCall.Method = method
	b.SendCall.BodyJSON, _ = json.Marshal(body)
	if len(result) > 0 {
		json.Unmarshal([]byte(b.SendCall.Result), result[0])
	}
	return b.SendCall.Err
}
