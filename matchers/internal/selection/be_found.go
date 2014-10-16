package selection

import (
	"fmt"
	"github.com/onsi/gomega/format"
)

type BeFoundMatcher struct{}

func (m *BeFoundMatcher) Match(actual interface{}) (success bool, err error) {
	actualSelection, ok := actual.(interface {
		Count() (int, error)
	})

	if !ok {
		return false, fmt.Errorf("BeFound matcher requires a Selection.  Got:\n%s", format.Object(actual, 1))
	}

	count, err := actualSelection.Count()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m *BeFoundMatcher) FailureMessage(actual interface{}) (message string) {
	return booleanSelectorMessage(actual, "to be found")
}

func (m *BeFoundMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return booleanSelectorMessage(actual, "not to be found")
}
