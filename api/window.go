package api

import "path"

type Window struct {
	ID      string
	Session *Session
}

func (w *Window) Send(endpoint, method string, body, result interface{}) error {
	return w.Session.Send(path.Join("window", w.ID, endpoint), method, body, result)
}

func (w *Window) SetSize(width, height int) error {
	request := struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}{width, height}

	return w.Send("size", "POST", request, nil)
}
