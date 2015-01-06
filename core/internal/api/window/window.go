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

func (w *Window) SwitchTo() error {
	request := struct {
		Handle string `json:"handle"`
	}{w.ID}
	if err := w.Session.Execute("window", "POST", &request, nil); err != nil {
		return err
	}
	return nil
}

func (w *Window) Close() error {
	err := w.SwitchTo()
	if err != nil {
		return err
	}

	if err := w.Session.Execute("window", "DELETE", nil, nil); err != nil {
		return err
	}
	return nil
}
