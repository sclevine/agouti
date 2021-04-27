package api

import "path"

type Window struct {
	ID      string
	Session *Session
}

func (w *Window) Send(method, endpoint string, body, result interface{}) error {
	return w.Session.Send(method, path.Join("window", w.ID, endpoint), body, result)
}

func (w *Window) SetSize(width, height int) error {
	request := struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}{width, height}

	return w.Send("POST", "size", request, nil)
}

func (w *Window) GetSize() (width int, height int, err error) {
	request := struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}{}


	if err := w.Send("GET", "size", nil, &request); err != nil {
		return 0, 0 , err
	}

	return request.Width, request.Height, nil
}
