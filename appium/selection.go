package appium

import (
	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/internal/target"
)

type Selection struct {
	*agouti.Selection
	session mobileSession
}

func (s *Selection) SelectionMethod() (string, error) {
	return "", nil
}

// Override Find to find by A11y instead of CSS selector
func (s *Selection) Find(id string) *Selection {
	return s.appendSelector(target.A11yID, id)
}

func (s *Selection) appendSelector(selectorType target.Type, value string) *Selection {
	selectors := target.Selectors(s.Selectors()).Append(selectorType, value)
	selection := s.WithSelectors(agouti.Selectors(selectors))
	return &Selection{selection, s.session}
}
