package selection

import (
	"fmt"

	"github.com/onsi/gomega/format"
)

type BeVisibleMatcher struct{}

func (m *BeVisibleMatcher) Match(actual interface{}) (success bool, err error) {
	actualSelection, ok := actual.(interface {
		Visible() (bool, error)
	})

	if !ok {
		return false, fmt.Errorf("BeVisible matcher requires a Selection.  Got:\n%s", format.Object(actual, 1))
	}

	visible, err := actualSelection.Visible()
	if err != nil {
		return false, err
	}

	return visible, nil
}

func (m *BeVisibleMatcher) FailureMessage(actual interface{}) (message string) {
	return booleanSelectorMessage(actual, "to be visible")
}

func (m *BeVisibleMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return booleanSelectorMessage(actual, "not to be visible")
}
