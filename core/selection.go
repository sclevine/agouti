package core

import "github.com/sclevine/agouti/core/internal/selection"

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
	Selectable

	// String returns a string representation of the selection, ex.
	//    CSS: .some-class | XPath: //table [3] | Link "click me" [single]
	String() string

	// Count returns the number of elements the selection refers to.
	Count() (int, error)

	// EqualsElement returns whether or not two selections of exactly
	// one element each refer to the same element.
	//	EqualsElement(comparable interface{}) (bool, error)

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

// A MultiSelection is a Selection with elements that may be indexed using
// the At(index) or Single() methods.
//
// A Selection returned by At(index) or Single() may still refer to multiple
// elements if any parent of the MultiSelection refers to multiple elements.
//
// Examples:
//    selection.All("section").All("form").At(1).Submit()
// Will submit the second form in each section.
//    selection.All("div").Find("h1").Click()
// Will click the only h1 in each div, failing if any div does not contain exactly one h1.
type MultiSelection interface {
	// All Selection methods are valid MultiSelection methods.
	Selection

	// At finds the element at the provided index.
	At(index int) Selection

	// Single specifies that the selection must refer to exactly one element.
	//    selection.Find("#selector")
	// is equivalent to
	//    selection.All("#selector").Single()
	Single() Selection
}

// Selectable methods apply their selectors to each element in the selection they
// are called on. If the selection they are called on refers to multiple
// elements, the resulting selection will refer to at least that many elements.
type Selectable interface {
	// For each element in the current selection:

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
}

type baseSelection struct {
	*selection.Selection
}

func (s *baseSelection) Find(selector string) Selection {
	return s.All(selector).Single()
}

func (s *baseSelection) FindByXPath(selector string) Selection {
	return s.AllByXPath(selector).Single()
}

func (s *baseSelection) FindByLink(text string) Selection {
	return s.AllByLink(text).Single()
}

func (s *baseSelection) FindByLabel(text string) Selection {
	return s.AllByLabel(text).Single()
}

func (s *baseSelection) First(selector string) Selection {
	return s.All(selector).At(0)
}

func (s *baseSelection) FirstByXPath(selector string) Selection {
	return s.AllByXPath(selector).At(0)
}

func (s *baseSelection) FirstByLink(text string) Selection {
	return s.AllByLink(text).At(0)
}

func (s *baseSelection) FirstByLabel(text string) Selection {
	return s.AllByLabel(text).At(0)
}

func (s *baseSelection) All(selector string) MultiSelection {
	return &multiSelection{&baseSelection{s.AppendCSS(selector)}}
}

func (s *baseSelection) AllByXPath(selector string) MultiSelection {
	return &multiSelection{&baseSelection{s.AppendXPath(selector)}}
}

func (s *baseSelection) AllByLink(text string) MultiSelection {
	return &multiSelection{&baseSelection{s.AppendLink(text)}}
}

func (s *baseSelection) AllByLabel(text string) MultiSelection {
	return &multiSelection{&baseSelection{s.AppendLabeled(text)}}
}

type multiSelection struct {
	*baseSelection
}

func (m *multiSelection) At(index int) Selection {
	return &baseSelection{m.AtIndex(index)}
}

func (m *multiSelection) Single() Selection {
	return &baseSelection{m.SingleOnly()}
}
