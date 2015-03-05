package appium

import "github.com/sclevine/agouti"

type Selection struct {
	*agouti.Selection
	session selectionSession
}

type selectionSession interface {
	SetEndpoint(thing string) error
}

func (s *Selection) SomeAppiumSelectionMethod() (string, error) {
	return "", nil
}

// override finder methods
func (s *Selection) Find(selector string) *Selection {
	return &Selection{s.Selection.Find(selector), s.session}
}
