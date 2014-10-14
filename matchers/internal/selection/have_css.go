package selection

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/sclevine/agouti/core"
)

type HaveCSSMatcher struct {
	ExpectedProperty string
	ExpectedValue    string
	actualValue      string
}

func (m *HaveCSSMatcher) Match(actual interface{}) (success bool, err error) {
	actualSelection, ok := actual.(core.Selection)
	if !ok {
		return false, fmt.Errorf("HaveCSS matcher requires a Selection.  Got:\n%s", format.Object(actual, 1))
	}

	m.actualValue, err = actualSelection.CSS(m.ExpectedProperty)
	if err != nil {
		return false, err
	}

	return m.actualValue == m.ExpectedValue, nil
}

func (m *HaveCSSMatcher) FailureMessage(actual interface{}) (message string) {
	return selectorMessage(actual, "to have CSS matching", m.style(m.ExpectedValue), m.style(m.actualValue))
}

func (m *HaveCSSMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return selectorMessage(actual, "not to have CSS matching", m.style(m.ExpectedValue), m.style(m.actualValue))
}

func (m *HaveCSSMatcher) style(value string) string {
	return fmt.Sprintf(`%s: "%s"`, m.ExpectedProperty, value)
}
