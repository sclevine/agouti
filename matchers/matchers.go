package matchers

import (
	"github.com/onsi/gomega/types"
	"github.com/sclevine/agouti/matchers/internal/selection"
)

func ContainText(text string) types.GomegaMatcher {
	return &selection.ContainTextMatcher{text}
}

func HaveAttribute(attribute string, value string) types.GomegaMatcher {
	return &selection.HaveAttributeMatcher{attribute, value}
}

func HaveCSS(property string, value string) types.GomegaMatcher {
	return &selection.HaveCSSMatcher{property, value}
}
