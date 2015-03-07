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

// Override Find to find by A11y
func (s *Selection) Find(id string) *Selection {
	selectors := target.Selectors(s.Selectors()).Append(target.A11yID, id)
	return &Selection{s.WithSelectors(agouti.Selectors(selectors)), s.session}
}
