package selection

import (
	"fmt"
	"github.com/onsi/gomega/format"
)

type HaveCSSMatcher struct {
	ExpectedProperty string
	ExpectedValue    string
	actualValue      string
}

type CSSer interface {
	CSS(property string) (string, error)
}

func (m *HaveCSSMatcher) Match(actual interface{}) (success bool, err error) {
	actualSelection, ok := actual.(CSSer)
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
