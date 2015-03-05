package internal

import (
	"fmt"

	"github.com/onsi/gomega/format"
	"github.com/sclevine/agouti"
)

type HaveLoggedErrorMatcher struct {
	ExpectedMessage string
}

func (m *HaveLoggedErrorMatcher) Match(actual interface{}) (success bool, err error) {
	actualPage, ok := actual.(interface {
		ReadAllLogs(logType string) ([]agouti.Log, error)
	})

	if !ok {
		return false, fmt.Errorf("HaveLoggedError matcher requires a Page.  Got:\n%s", format.Object(actual, 1))
	}

	logs, err := actualPage.ReadAllLogs("browser")
	if err != nil {
		return false, err
	}

	matchAnyMessage := m.ExpectedMessage == ""

	for _, log := range logs {
		levelMatches := log.Level == "WARNING" || log.Level == "SEVERE"
		logMatches := matchAnyMessage || log.Message == m.ExpectedMessage
		if levelMatches && logMatches {
			return true, nil
		}
	}

	return false, nil
}

func (m *HaveLoggedErrorMatcher) FailureMessage(actual interface{}) (message string) {
	return equalityMessage(actual, "to have error log matching", m.ExpectedMessage)

}

func (m *HaveLoggedErrorMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return equalityMessage(actual, "not to have error log matching", m.ExpectedMessage)
}
