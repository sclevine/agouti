package mocks

type Element struct {
	GetTextCall struct {
		ReturnText string
		Err error
	}
}

func (e *Element) GetText() (string, error) {
	return e.GetTextCall.ReturnText, e.GetTextCall.Err
}
