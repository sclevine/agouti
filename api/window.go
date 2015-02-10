package api

type Window struct {
	ID      string
	Session *Session
}

func (w *Window) SetSize(width, height int) error {
	request := struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}{width, height}

	return w.Session.sendWindow(w.ID, "size", "POST", request)
}
