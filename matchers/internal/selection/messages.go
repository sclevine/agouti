package selection

import (
	"fmt"
	"github.com/onsi/gomega/format"
)

func selectorMessage(actual interface{}, message, expected, actualValue string) string {
	failureMessage := "Expected selection '%s' %s\n%s%s\nbut found\n%s%s"
	return fmt.Sprintf(failureMessage, actual, message, format.Indent, expected, format.Indent, actualValue)
}

func binarySelectorMessage(actual interface{}, message string, expected interface{}) string {
	failureMessage := "Expected selection '%s' %s\n%s%s"
	return fmt.Sprintf(failureMessage, actual, message, format.Indent, expected)
}

func booleanSelectorMessage(actual interface{}, message string) string {
	failureMessage := "Expected selection '%s' %s"
	return fmt.Sprintf(failureMessage, actual, message)
}
