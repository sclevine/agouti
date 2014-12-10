package page

import (
	"fmt"
	"github.com/onsi/gomega/format"
)

func pageMessage(message, expected string, actualValue ...string) string {
	if len(actualValue) == 0 {
		failureMessage := "Expected page %s\n%s%s"
		return fmt.Sprintf(failureMessage, message, format.Indent, expected)
	} else {
		failureMessage := "Expected page %s\n%s%s\nbut found\n%s%s"
		return fmt.Sprintf(failureMessage, message, format.Indent, expected, format.Indent, actualValue[0])
	}
}
