package window

type Executable interface {
	Execute(endpoint, method string, body, result interface{}) error
}

type Window struct {
	ID string
	Session Executable
}

func (w *Window) SetSize(height, width int) error {
	endpoint := "window/" + w.ID + "/size"
	request := struct {
			Width int `json:"width"`
			Height int `json:"height"`
		}{width,height}

	if err := w.Session.Execute(endpoint, "POST", &request, &struct{}{}); err != nil {
		return err
	}
	return nil
}