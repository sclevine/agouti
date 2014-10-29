package core

import (
	"fmt"

	"github.com/sclevine/agouti/core/internal/types"
)

// Selection instances refer to a selection of elements.
// All Selection methods are valid MultiSelection methods.
//
// Examples:
//
//    selection.Find("table").All("tr").At(2).First("td input[type=checkbox]").Check()
// Will check the first checkbox in the third row of the only table.
//    selection.Find("table").All("tr").Find("td").All("input[type=checkbox]").Check()
// Will check all checkboxes in the first-and-only cell of each row in the only table.
type Selection interface {
	// The Find<X>(), First<X>(), and All<X>() methods apply their selectors to
	// each element in the selection they are called on. If the selection they are
	// called on refers to multiple elements, the resulting selection will refer
	// to at least that many elements.

	// Therefore, for each element in the current selection:

	// Find finds exactly one element by CSS selector.
	Find(selector string) Selection

	// FindByXPath finds exactly one element by XPath selector.
	FindByXPath(selector string) Selection

	// FindByLink finds exactly one element by anchor link text.
	FindByLink(text string) Selection

	// FindByLabel finds exactly one element by associated label text.
	FindByLabel(text string) Selection

	// First finds the first element by CSS selector.
	First(selector string) Selection

	// FirstByXPath finds the first element by XPath selector.
	FirstByXPath(selector string) Selection

	// FirstByLink finds the first element by anchor link text.
	FirstByLink(text string) Selection

	// FirstByLabel finds the first element by associated label text.
	FirstByLabel(text string) Selection

	// All finds zero or more elements by CSS selector.
	All(selector string) MultiSelection

	// AllByXPath finds zero or more elements by XPath selector.
	AllByXPath(selector string) MultiSelection

	// AllByLink finds zero or more elements by anchor link text.
	AllByLink(text string) MultiSelection

	// AllByLabel finds zero or more elements by associated label text.
	AllByLabel(text string) MultiSelection

	// String returns a string representation of the selection, ex.
	//    CSS: .some-class | XPath: //table [3] | Link "click me" [single]
	String() string

	// Count returns the number of elements the selection refers to.
	Count() (int, error)

	// EqualsElement returns whether or not two selections of exactly
	// one element each refer to the same element.
	EqualsElement(comparable interface{}) (bool, error)

	// Click clicks on all of the elements the selection refers to.
	Click() error

	// DoubleClick double-clicks on all of the elements the selection refers to.
	DoubleClick() error

	// Fill fills all of the input fields the selection refers to.
	Fill(text string) error

	// Check checks all of the unchecked checkboxes that the selection refers to.
	Check() error

	// Uncheck unchecks all of the checked checkboxes that the selection refers to.
	Uncheck() error

	// Select, when called on any number of <select> elements, will select all
	// <options> under those elements that match the provided text.
	Select(text string) error

	// Submit submits a form. The selection may refer to the form itself, or any
	// input element contained within the form.
	Submit() error

	// Text returns text for exactly one element.
	Text() (string, error)

	// Attribute returns an attribute value for exactly one element.
	Attribute(attribute string) (string, error)

	// CSS returns a CSS style property value for exactly one element.
	CSS(property string) (string, error)

	// Selected returns true if all of the elements that the selection
	// refers to are selected.
	Selected() (bool, error)

	// Visible returns true if all of the elements that the selection
	// refers to are visible.
	Visible() (bool, error)

	// Enabled returns true if all of the elements that the selection
	// refers to are enabled.
	Enabled() (bool, error)
}

type selection struct {
	client    types.Client
	selectors []types.Selector
}

func (s *selection) Find(selector string) Selection {
	return s.All(selector).Single()
}

func (s *selection) FindByXPath(selector string) Selection {
	return s.AllByXPath(selector).Single()
}

func (s *selection) FindByLink(text string) Selection {
	return s.AllByLink(text).Single()
}

func (s *selection) FindByLabel(text string) Selection {
	return s.AllByLabel(text).Single()
}

func (s *selection) First(selector string) Selection {
	return s.All(selector).At(0)
}

func (s *selection) FirstByXPath(selector string) Selection {
	return s.AllByXPath(selector).At(0)
}

func (s *selection) FirstByLink(text string) Selection {
	return s.AllByLink(text).At(0)
}

func (s *selection) FirstByLabel(text string) Selection {
	return s.AllByLabel(text).At(0)
}

func (s *selection) All(selector string) MultiSelection {
	last := len(s.selectors) - 1

	lastIsCSS := last >= 0 && s.selectors[last].Using == "css selector"
	if lastIsCSS && !s.selectors[last].Indexed && !s.selectors[last].Single {
		return s.mergedSelection(selector)
	}

	return s.subSelection("css selector", selector)
}

func (s *selection) AllByXPath(selector string) MultiSelection {
	return s.subSelection("xpath", selector)
}

func (s *selection) AllByLink(text string) MultiSelection {
	return s.subSelection("link text", text)
}

func (s *selection) AllByLabel(text string) MultiSelection {
	selector := fmt.Sprintf(`//input[@id=(//label[normalize-space(text())="%s"]/@for)] | //label[normalize-space(text())="%s"]/input`, text, text)
	return s.AllByXPath(selector)
}

func (s *selection) subSelection(using, value string) MultiSelection {
	newSelector := types.Selector{Using: using, Value: value}
	selection := &selection{s.client, appendSelector(s.selectors, newSelector)}
	return &multiSelection{selection}
}

func (s *selection) mergedSelection(value string) MultiSelection {
	last := len(s.selectors) - 1
	newSelectorValue := s.selectors[last].Value + " " + value
	newSelector := types.Selector{Using: "css selector", Value: newSelectorValue}
	selection := &selection{s.client, appendSelector(s.selectors[:last], newSelector)}
	return &multiSelection{selection}
}

func appendSelector(selectors []types.Selector, selector types.Selector) []types.Selector {
	selectorsCopy := append([]types.Selector(nil), selectors...)
	return append(selectorsCopy, selector)
}
