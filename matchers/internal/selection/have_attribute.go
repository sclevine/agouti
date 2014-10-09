package selection

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/sclevine/agouti/page"
)

type HaveAttributeMatcher struct {
	ExpectedAttribute string
	ExpectedValue     string
}

func (m *HaveAttributeMatcher) Match(actual interface{}) (success bool, err error) {
	actualPage, ok := actual.(page.Selection)
	if !ok {
		return false, fmt.Errorf("HaveAttribute matcher requires a Selection or Page.  Got:\n%s", format.Object(actual, 1))
	}

	actualValue, err := actualPage.Attribute(m.ExpectedAttribute)
	if err != nil {
		return false, err
	}

	return actualValue == m.ExpectedValue, nil
}

func (m *HaveAttributeMatcher) FailureMessage(actual interface{}) (message string) {
	actualSelector := SelectorText(actual.(page.Selection).Selector())
	return format.Message(actualSelector, "to have attribute matching", m.attributeSelector())
}

func (m *HaveAttributeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	actualSelector := SelectorText(actual.(page.Selection).Selector())
	return format.Message(actualSelector, "not to have attribute matching", m.attributeSelector())
}

func (m *HaveAttributeMatcher) attributeSelector() string {
	return fmt.Sprintf(`[%s="%s"]`, m.ExpectedAttribute, m.ExpectedValue)
}
