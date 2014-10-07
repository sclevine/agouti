package mocks

type Element struct {
	GetTextCall struct {
		ReturnText string
		Err        error
	}

	ClickCall struct {
		Called bool
		Err    error
	}
}

func (e *Element) GetText() (string, error) {
	return e.GetTextCall.ReturnText, e.GetTextCall.Err
}

func (e *Element) Click() error {
	e.ClickCall.Called = true
	return e.ClickCall.Err
}
