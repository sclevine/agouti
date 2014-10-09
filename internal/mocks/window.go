package mocks

type Window struct {
	SizeCall struct {
		Width  int
		Height int
		Err    error
	}
}

func (w *Window) SetSize(height, width int) error {
	w.SizeCall.Height = height
	w.SizeCall.Width = width
	return w.SizeCall.Err
}
