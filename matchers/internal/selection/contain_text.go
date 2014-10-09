package selection

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/sclevine/agouti/page"
)

type ContainTextMatcher struct {
	ExpectedText string
}

func (m *ContainTextMatcher) Match(actual interface{}) (success bool, err error) {
	actualPage, ok := actual.(page.Selection)
	if !ok {
		return false, fmt.Errorf("ContainText matcher requires a Selection or Page.  Got:\n%s", format.Object(actual, 1))
	}

	actualText, err := actualPage.Text()
	if err != nil {
		return false, err
	}

	return actualText == m.ExpectedText, nil
}

func (m *ContainTextMatcher) FailureMessage(actual interface{}) (message string) {
	actualSelector := SelectorText(actual.(page.Selection).Selector())
	return format.Message(actualSelector, "to have text matching", m.ExpectedText)
}

func (m *ContainTextMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	actualSelector := SelectorText(actual.(page.Selection).Selector())
	return format.Message(actualSelector, "not to have text matching", m.ExpectedText)
}
