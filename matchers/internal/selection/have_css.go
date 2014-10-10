package selection

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/sclevine/agouti/page"
)

type HaveCSSMatcher struct {
	ExpectedProperty string
	ExpectedValue    string
}

func (m *HaveCSSMatcher) Match(actual interface{}) (success bool, err error) {
	actualPage, ok := actual.(page.Selection)
	if !ok {
		return false, fmt.Errorf("HaveCSS matcher requires a Selection or Page.  Got:\n%s", format.Object(actual, 1))
	}

	actualValue, err := actualPage.CSS(m.ExpectedProperty)
	if err != nil {
		return false, err
	}

	return actualValue == m.ExpectedValue, nil
}

func (m *HaveCSSMatcher) FailureMessage(actual interface{}) (message string) {
	actualSelector := SelectorText(actual.(page.Selection).Selector())
	return format.Message(actualSelector, "to have CSS matching", m.style())
}

func (m *HaveCSSMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	actualSelector := SelectorText(actual.(page.Selection).Selector())
	return format.Message(actualSelector, "not to have CSS matching", m.style())
}

func (m *HaveCSSMatcher) style() string {
	return fmt.Sprintf(`%s: "%s"`, m.ExpectedProperty, m.ExpectedValue)
}
