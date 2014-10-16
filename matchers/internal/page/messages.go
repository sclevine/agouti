package page

import (
	"fmt"
	"github.com/onsi/gomega/format"
)

func pageMessage(message, expected, actualValue string) string {
	failureMessage := "Expected page %s\n%s%s\nbut found\n%s%s"
	return fmt.Sprintf(failureMessage, message, format.Indent, expected, format.Indent, actualValue)
}
