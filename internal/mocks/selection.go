package mocks

import "github.com/sclevine/agouti/page"

type Selection struct {
	SelectorCall struct {
		ReturnSelector string
	}

	TextCall struct {
		ReturnText string
		Err        error
	}

	AttributeCall struct {
		Attribute   string
		ReturnValue string
		Err         error
	}

	CSSCall struct {
		Property    string
		ReturnValue string
		Err         error
	}

	SelectedCall struct {
		ReturnSelected bool
		Err            error
	}
}

func (s *Selection) Find(selector string) page.Selection {
	return &Selection{}
}

func (s *Selection) Selector() string {
	return s.SelectorCall.ReturnSelector
}

func (s *Selection) Click() error {
	return nil
}

func (s *Selection) Check() error {
	return nil
}

func (s *Selection) Fill(text string) error {
	return nil
}

func (s *Selection) Text() (string, error) {
	return s.TextCall.ReturnText, s.TextCall.Err
}

func (s *Selection) Attribute(attribute string) (string, error) {
	s.AttributeCall.Attribute = attribute
	return s.AttributeCall.ReturnValue, s.AttributeCall.Err
}

func (s *Selection) CSS(property string) (string, error) {
	s.CSSCall.Property = property
	return s.CSSCall.ReturnValue, s.CSSCall.Err
}

func (s *Selection) Selected() (bool, error) {
	return s.SelectedCall.ReturnSelected, s.SelectedCall.Err
}
