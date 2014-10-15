package mocks

type Window struct {
	SizeCall struct {
		Width  int
		Height int
		Err    error
	}
}

func (w *Window) SetSize(width, height int) error {
	w.SizeCall.Width = width
	w.SizeCall.Height = height
	return w.SizeCall.Err
}
