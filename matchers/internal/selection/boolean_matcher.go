package selection

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/onsi/gomega/format"
)

type BooleanMatcher struct {
	Method string
}

func (m *BooleanMatcher) Match(actual interface{}) (success bool, err error) {
	method := reflect.ValueOf(actual).MethodByName(m.Method)
	if !method.IsValid() {
		return false, fmt.Errorf("Matcher requires a *Selection.  Got:\n%s", format.Object(actual, 1))
	}

	results := method.Call(nil)
	matchValue := results[0]
	errValue := results[1]

	if !errValue.IsNil() {
		return false, errValue.Interface().(error)
	}

	return matchValue.Bool(), nil
}

func (m *BooleanMatcher) FailureMessage(actual interface{}) (message string) {
	return booleanSelectorMessage(actual, "to be "+strings.ToLower(m.Method))
}

func (m *BooleanMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return booleanSelectorMessage(actual, "not to be "+strings.ToLower(m.Method))
}
