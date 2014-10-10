package selection

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/sclevine/agouti/page"
)

func selectorMessage(actual interface{}, message, expected, actualValue string) string {
	actualSelector := actual.(page.Selection).Selector()
	failureMessage := "Expected selection '%s' %s\n%s%s\nbut found\n%s%s"
	return fmt.Sprintf(failureMessage, actualSelector, message, format.Indent, expected, format.Indent, actualValue)
}

func booleanSelectorMessage(actual interface{}, message string) string {
	actualSelector := actual.(page.Selection).Selector()
	failureMessage := "Expected selection '%s' %s"
	return fmt.Sprintf(failureMessage, actualSelector, message)
}
