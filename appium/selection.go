package appium

import (
	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api/mobile"
	"github.com/sclevine/agouti/internal/target"
)

type Selection struct {
	*agouti.Selection
	session *mobile.Session
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

func (s *Selection) FindByID(id string) *Selection {
	sel := target.Selectors{}
	apiSession := s.session.Session
	return &Selection{agouti.NewSelection(apiSession, sel), s.session}
}
