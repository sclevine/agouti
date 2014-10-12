package mocks

type Element struct {
	GetTextCall struct {
		ReturnText string
		Err        error
	}

	GetAttributeCall struct {
		Attribute   string
		ReturnValue string
		Err         error
	}

	GetCSSCall struct {
		Property    string
		ReturnValue string
		Err         error
	}

	ClickCall struct {
		Called bool
		Err    error
	}

	ClearCall struct {
		Called bool
		Err    error
	}

	ValueCall struct {
		Text string
		Err  error
	}

	IsSelectedCall struct {
		ReturnSelected bool
		Err            error
	}

	SubmitCall struct {
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

func (e *Element) GetCSS(property string) (string, error) {
	e.GetCSSCall.Property = property
	return e.GetCSSCall.ReturnValue, e.GetCSSCall.Err
}

func (e *Element) Click() error {
	e.ClickCall.Called = true
	return e.ClickCall.Err
}

func (e *Element) Clear() error {
	e.ClearCall.Called = true
	return e.ClearCall.Err
}

func (e *Element) Value(text string) error {
	e.ValueCall.Text = text
	return e.ValueCall.Err
}

func (e *Element) IsSelected() (bool, error) {
	return e.IsSelectedCall.ReturnSelected, e.IsSelectedCall.Err
}

func (e *Element) Submit() error {
	e.SubmitCall.Called = true
	return e.SubmitCall.Err
}
