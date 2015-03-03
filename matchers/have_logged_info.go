package matchers

import (
	"fmt"

	"github.com/onsi/gomega/format"
	"github.com/sclevine/agouti"
)

type HaveLoggedInfoMatcher struct {
	ExpectedMessage string
}

func (m *HaveLoggedInfoMatcher) Match(actual interface{}) (success bool, err error) {
	actualPage, ok := actual.(interface {
		ReadAllLogs(logType string) ([]agouti.Log, error)
	})

	if !ok {
		return false, fmt.Errorf("HaveLoggedInfo matcher requires a Page.  Got:\n%s", format.Object(actual, 1))
	}

	logs, err := actualPage.ReadAllLogs("browser")
	if err != nil {
		return false, err
	}

	matchAnyMessage := m.ExpectedMessage == ""

	for _, log := range logs {
		logMatches := matchAnyMessage || log.Message == m.ExpectedMessage
		if log.Level == "INFO" && logMatches {
			return true, nil
		}
	}

	return false, nil
}

func (m *HaveLoggedInfoMatcher) FailureMessage(actual interface{}) (message string) {
	return equalityMessage(actual, "to have info log matching", m.ExpectedMessage)

}

func (m *HaveLoggedInfoMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return equalityMessage(actual, "not to have info log matching", m.ExpectedMessage)
}
