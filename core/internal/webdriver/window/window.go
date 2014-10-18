package window

type Window struct {
	ID      string
	Session session
}

type session interface {
	Execute(endpoint, method string, body interface{}, result ...interface{}) error
}

func (w *Window) SetSize(width, height int) error {
	endpoint := "window/" + w.ID + "/size"
	request := struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}{width, height}

	if err := w.Session.Execute(endpoint, "POST", &request, &struct{}{}); err != nil {
		return err
	}
	return nil
}
