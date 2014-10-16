package selection

import (
	"fmt"
	"github.com/onsi/gomega/format"
)

type BeSelectedMatcher struct{}

type Selectable interface {
	Selected() (bool, error)
}

func (m *BeSelectedMatcher) Match(actual interface{}) (success bool, err error) {
	actualPage, ok := actual.(Selectable)
	if !ok {
		return false, fmt.Errorf("BeSelected matcher requires a Selection.  Got:\n%s", format.Object(actual, 1))
	}

	selected, err := actualPage.Selected()
	if err != nil {
		return false, err
	}

	return selected, nil
}

func (m *BeSelectedMatcher) FailureMessage(actual interface{}) (message string) {
	return booleanSelectorMessage(actual, "to be selected")
}

func (m *BeSelectedMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return booleanSelectorMessage(actual, "not to be selected")
}
