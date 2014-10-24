package page

import (
	"fmt"

	"github.com/onsi/gomega/format"
)

type HaveTitleMatcher struct {
	ExpectedTitle string
	actualTitle   string
}

func (m *HaveTitleMatcher) Match(actual interface{}) (success bool, err error) {
	actualPage, ok := actual.(interface {
		Title() (string, error)
	})

	if !ok {
		return false, fmt.Errorf("HaveTitle matcher requires a Page.  Got:\n%s", format.Object(actual, 1))
	}

	m.actualTitle, err = actualPage.Title()
	if err != nil {
		return false, err
	}

	return m.actualTitle == m.ExpectedTitle, nil
}

func (m *HaveTitleMatcher) FailureMessage(_ interface{}) (message string) {
	return pageMessage("to have title matching", m.ExpectedTitle, m.actualTitle)

}

func (m *HaveTitleMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return pageMessage("not to have title matching", m.ExpectedTitle, m.actualTitle)
}
