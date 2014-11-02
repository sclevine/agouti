package page

import (
	"fmt"

	"github.com/onsi/gomega/format"
)

type HaveURLMatcher struct {
	ExpectedURL string
	actualURL   string
}

func (m *HaveURLMatcher) Match(actual interface{}) (success bool, err error) {
	actualPage, ok := actual.(interface {
		URL() (string, error)
	})

	if !ok {
		return false, fmt.Errorf("HaveURL matcher requires a Page.  Got:\n%s", format.Object(actual, 1))
	}

	m.actualURL, err = actualPage.URL()
	if err != nil {
		return false, err
	}

	return m.actualURL == m.ExpectedURL, nil
}

func (m *HaveURLMatcher) FailureMessage(_ interface{}) (message string) {
	return pageMessage("to have URL matching", m.ExpectedURL, m.actualURL)

}

func (m *HaveURLMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return pageMessage("not to have URL matching", m.ExpectedURL, m.actualURL)
}
