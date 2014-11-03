package matchers

import (
	"github.com/onsi/gomega/types"
	"github.com/sclevine/agouti/matchers/internal/selection"
)

// HaveText passes when the expected text is equal to the actual element text.
// This matcher will fail if the provided selection refers to more than one element.
func HaveText(text string) types.GomegaMatcher {
	return &selection.HaveTextMatcher{ExpectedText: text}
}

// MatchText passes when the expected regular expression matches the actual element text.
// This matcher will fail if the provided selection refers to more than one element.
func MatchText(regexp string) types.GomegaMatcher {
	return &selection.MatchTextMatcher{Regexp: regexp}
}

// HaveAttribute passes when the expected attribute and value are present on the element.
// This matcher will fail if the provided selection refers to more than one element.
func HaveAttribute(attribute string, value string) types.GomegaMatcher {
	return &selection.HaveAttributeMatcher{ExpectedAttribute: attribute, ExpectedValue: value}
}

// HaveCSS passes when the expected CSS property and value are present on the element.
// This matcher only matches exact, calculated CSS values, though there is support for parsing colors.
// Example: "blue" and "#00f" will both match rgba(0, 0, 255, 1)
// This matcher will fail if the provided selection refers to more than one element.
func HaveCSS(property string, value string) types.GomegaMatcher {
	return &selection.HaveCSSMatcher{ExpectedProperty: property, ExpectedValue: value}
}

// BeSelected passes when the provided selection refers to form elements that are selected.
// Examples: a checked <input type="checkbox" />, or the selected <option> in a <select>
// This matcher will fail if any of the selection's form elements are not selected.
func BeSelected() types.GomegaMatcher {
	return &selection.BeSelectedMatcher{}
}

// BeVisible passes when the selection refers to elements that are displayed on the page.
// This matcher will fail if any of the selection's elements are not visible.
func BeVisible() types.GomegaMatcher {
	return &selection.BeVisibleMatcher{}
}

// BeEnabled passes when the selection refers to form elements that are enabled.
// This matcher will fail if any of the selection's form elements are not enabled.
func BeEnabled() types.GomegaMatcher {
	return &selection.BeEnabledMatcher{}
}

// BeActive passes when the selection refers to the active page element.
func BeActive() types.GomegaMatcher {
	return &selection.BeActiveMatcher{}
}

// BeFound passes when the provided selection refers to one or more elements on the page.
func BeFound() types.GomegaMatcher {
	return &selection.BeFoundMatcher{}
}

// EqualElement passes when the expected selection refers to the same element as the provided
// actual selection. This matcher will fail if either selection refers to more than one element.
func EqualElement(comparable interface{}) types.GomegaMatcher {
	return &selection.EqualElementMatcher{ExpectedSelection: comparable}
}
