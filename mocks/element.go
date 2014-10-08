package mocks

type Element struct {
	GetTextCall struct {
		ReturnText string
		Err        error
	}

	GetAttributeCall struct {
		Attribute    string
		ReturnValue     string
		Err        		error
	}

	ClickCall struct {
		Called bool
		Err    error
	}
}

func (e *Element) GetText() (string, error) {
	return e.GetTextCall.ReturnText, e.GetTextCall.Err
}

func (e *Element) GetAttribute(attribute string) (string, error) {
	e.GetAttributeCall.Attribute = attribute
	return e.GetAttributeCall.ReturnValue, e.GetAttributeCall.Err
}

func (e *Element) Click() error {
	e.ClickCall.Called = true
	return e.ClickCall.Err
}
