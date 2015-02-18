package page

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

func (m *HaveLoggedInfoMatcher) FailureMessage(_ interface{}) (message string) {
	return pageMessage("to have info log matching", m.ExpectedMessage)

}

func (m *HaveLoggedInfoMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return pageMessage("not to have info log matching", m.ExpectedMessage)
}
