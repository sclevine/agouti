package selection

import (
	"fmt"

	"github.com/onsi/gomega/format"
)

type BeEnabledMatcher struct{}

func (m *BeEnabledMatcher) Match(actual interface{}) (success bool, err error) {
	actualSelection, ok := actual.(interface {
		Enabled() (bool, error)
	})

	if !ok {
		return false, fmt.Errorf("BeEnabled matcher requires a Selection.  Got:\n%s", format.Object(actual, 1))
	}

	enabled, err := actualSelection.Enabled()
	if err != nil {
		return false, err
	}

	return enabled, nil
}

func (m *BeEnabledMatcher) FailureMessage(actual interface{}) (message string) {
	return booleanSelectorMessage(actual, "to be enabled")
}

func (m *BeEnabledMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return booleanSelectorMessage(actual, "not to be enabled")
}
