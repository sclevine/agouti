package selection

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/sclevine/agouti/page"
)

type ContainTextMatcher struct {
	ExpectedText string
	actualText   string
}

func (m *ContainTextMatcher) Match(actual interface{}) (success bool, err error) {
	actualPage, ok := actual.(page.Selection)
	if !ok {
		return false, fmt.Errorf("ContainText matcher requires a Selection or Page.  Got:\n%s", format.Object(actual, 1))
	}

	m.actualText, err = actualPage.Text()
	if err != nil {
		return false, err
	}

	return m.actualText == m.ExpectedText, nil
}

func (m *ContainTextMatcher) FailureMessage(actual interface{}) (message string) {
	return selectorMessage(actual, "to have text matching", m.ExpectedText, m.actualText)

}

func (m *ContainTextMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return selectorMessage(actual, "not to have text matching", m.ExpectedText, m.actualText)
}
