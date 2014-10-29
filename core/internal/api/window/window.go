package window

import "github.com/sclevine/agouti/core/internal/types"

type Window struct {
	ID      string
	Session types.Session
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
