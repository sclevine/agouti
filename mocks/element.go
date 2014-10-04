package mocks

type Element struct {
	GetText struct {
		ReturnText string
		Err error
	}
}

func (e *Element) GetText() (string, error) {
	return e.ReturnText, e.Err
}
