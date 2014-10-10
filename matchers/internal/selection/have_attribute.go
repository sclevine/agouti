package selection

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/sclevine/agouti/page"
)

type HaveAttributeMatcher struct {
	ExpectedAttribute string
	ExpectedValue     string
	actualValue       string
}

func (m *HaveAttributeMatcher) Match(actual interface{}) (success bool, err error) {
	actualPage, ok := actual.(page.Selection)
	if !ok {
		return false, fmt.Errorf("HaveAttribute matcher requires a Selection or Page.  Got:\n%s", format.Object(actual, 1))
	}

	m.actualValue, err = actualPage.Attribute(m.ExpectedAttribute)
	if err != nil {
		return false, err
	}

	return m.actualValue == m.ExpectedValue, nil
}

func (m *HaveAttributeMatcher) FailureMessage(actual interface{}) (message string) {
	return selectorMessage(actual, "to have attribute matching", m.attribute(m.ExpectedValue), m.attribute(m.actualValue))
}

func (m *HaveAttributeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return selectorMessage(actual, "not to have attribute matching", m.attribute(m.ExpectedValue), m.attribute(m.actualValue))
}

func (m *HaveAttributeMatcher) attribute(value string) string {
	return fmt.Sprintf(`[%s="%s"]`, m.ExpectedAttribute, value)
}
