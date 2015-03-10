package appium

import (
	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/internal/target"
)

type Selection struct {
	*agouti.Selection
	session mobileSession
}

// Agouti wrapped selectors

func (s *Selection) Find(selector string) *Selection {
	return s.wrap(s.Selection.Find(selector))
}

func (s *Selection) FindByID(selector string) *Selection {
	return s.wrap(s.Selection.FindByID(selector))
}

func (s *Selection) FindByXPath(xPath string) *Selection {
	return s.wrap(s.Selection.FindByXPath(xPath))
}

func (s *Selection) FindByLink(text string) *Selection {
	return s.wrap(s.Selection.FindByLink(text))
}

// Appium specific selectors

func (s *Selection) FindByA11yID(id string) *Selection {
	return s.newSelection(s.appendSelector(target.A11yID, id).Single())
}

func (s *Selection) FindByAndroidUI(uiautomatorQuery string) *Selection {
	return s.newSelection(s.appendSelector(target.AndroidAut, uiautomatorQuery).Single())
}

func (s *Selection) FindByiOSUI(uiautomationQuery string) *Selection {
	return s.newSelection(s.appendSelector(target.IOSAut, uiautomationQuery).Single())
}



// Selection helpers
func (s *Selection) appendSelector(selectorType target.Type, value string) target.Selectors {
	return target.Selectors(s.Selectors()).Append(selectorType, value)
}
func (s *Selection) newSelection(selectors target.Selectors) *Selection {
	return &Selection{s.WithSelectors(agouti.Selectors(selectors)), s.session}
}
func (s *Selection) wrap(selection *agouti.Selection) *Selection {
	return &Selection{selection, s.session}
}
