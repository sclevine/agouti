package selection

import (
	"fmt"

	"github.com/onsi/gomega/format"
)

type BeActiveMatcher struct{}

func (m *BeActiveMatcher) Match(actual interface{}) (success bool, err error) {
	actualSelection, ok := actual.(interface {
		Active() (bool, error)
	})

	if !ok {
		return false, fmt.Errorf("BeActive matcher requires a Selection.  Got:\n%s", format.Object(actual, 1))
	}

	active, err := actualSelection.Active()
	if err != nil {
		return false, err
	}

	return active, nil
}

func (m *BeActiveMatcher) FailureMessage(actual interface{}) (message string) {
	return booleanSelectorMessage(actual, "to be active")
}

func (m *BeActiveMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return booleanSelectorMessage(actual, "not to be active")
}
