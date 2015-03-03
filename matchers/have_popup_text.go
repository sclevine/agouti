package matchers

import (
	"fmt"

	"github.com/onsi/gomega/format"
)

type HavePopupTextMatcher struct {
	ExpectedText string
	actualText   string
}

func (m *HavePopupTextMatcher) Match(actual interface{}) (success bool, err error) {
	actualPage, ok := actual.(interface {
		PopupText() (string, error)
	})

	if !ok {
		return false, fmt.Errorf("HavePopupText matcher requires a Page.  Got:\n%s", format.Object(actual, 1))
	}

	m.actualText, err = actualPage.PopupText()
	if err != nil {
		return false, err
	}

	return m.actualText == m.ExpectedText, nil
}

func (m *HavePopupTextMatcher) FailureMessage(actual interface{}) (message string) {
	return valueMessage(actual, "to have popup text matching", m.ExpectedText, m.actualText)

}

func (m *HavePopupTextMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return valueMessage(actual, "not to have popup text matching", m.ExpectedText, m.actualText)
}
