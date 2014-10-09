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

func (s *Selection) Text() (string, error) {
	return s.TextCall.ReturnText, s.TextCall.Err
}

func (s *Selection) Attribute(attribute string) (string, error) {
	s.AttributeCall.Attribute = attribute
	return s.AttributeCall.ReturnValue, s.AttributeCall.Err
}
