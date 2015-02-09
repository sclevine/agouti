package core

// Selection instances refer to a selection of elements.
// Every Selection method is a valid MultiSelection method.
//
// Examples:
//
//    selection.Find("table").All("tr").At(2).First("td input[type=checkbox]").Check()
// Checks the first checkbox in the third row of the only table.
//    selection.Find("table").All("tr").Find("td").All("input[type=checkbox]").Check()
// Checks all checkboxes in the first-and-only cell of each row in the only table.
type Selection interface {
	// Selections have Find*(), All*(), and First*() methods that return new sub-Selections
	Selectable

	// String returns a string representation of the selection, ex.
	//    CSS: .some-class | XPath: //table [3] | Link "click me" [single]
	String() string

	// Count returns the number of elements that the selection refers to.
	Count() (int, error)

	// EqualsElement returns whether or not two selections of exactly
	// one element each refer to the same element.
	EqualsElement(comparable interface{}) (bool, error)

	// SwitchToFrame focuses on the frame specified by the selection. All new and
	// existing selections will refer to the new frame. All further Page methods
	// will apply to this frame as well.
	SwitchToFrame() error

	// Click clicks on all of the elements that the selection refers to.
	Click() error

	// DoubleClick double-clicks on all of the elements that the selection refers to.
	DoubleClick() error

	// Fill fills all of the fields the selection refers to with the provided text.
	Fill(text string) error

	// Check checks all of the unchecked checkboxes that the selection refers to.
	Check() error

	// Uncheck unchecks all of the checked checkboxes that the selection refers to.
	Uncheck() error

	// Select, when called on some number of <select> elements, will select all
	// <option> elements under those <select> elements that match the provided text.
	Select(text string) error

	// Submit submits all selected forms. The selection may refer to a form itself
	// or any input element contained within a form.
	Submit() error

	// Text returns the entirety of the text content for exactly one element.
	Text() (string, error)

	// Attribute returns an attribute value for exactly one element.
	Attribute(attribute string) (string, error)

	// CSS returns a CSS style property value for exactly one element.
	CSS(property string) (string, error)

	// Selected returns true if all of the elements that the selection refers to
	// are selected.
	Selected() (bool, error)

	// Visible returns true if all of the elements that the selection refers to
	// are visible.
	Visible() (bool, error)

	// Enabled returns true if all of the elements that the selection refers to
	// are enabled.
	Enabled() (bool, error)

	// Active returns true if the single element that the selection refers to is active.
	Active() (bool, error)
}

// A MultiSelection is a Selection that may be indexed using the At() method.
//
// A Selection returned by At() may still refer to multiple elements if any
// parent of the MultiSelection refers to multiple elements.
//
// Examples:
//    selection.All("section").All("form").At(1).Submit()
// Submits the second form in each section.
//    selection.All("div").Find("h1").Click()
// Clicks one h1 in each div, failing if any div does not contain exactly one h1.
type MultiSelection interface {
	// All Selection methods are valid MultiSelection methods.
	Selection

	// At finds the element at the provided index.
	At(index int) Selection
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

	// FindByLink finds exactly one anchor element by its text content.
	FindByLink(text string) Selection

	// FindByLabel finds exactly one element by associated label text.
	FindByLabel(text string) Selection

	// FindByButton finds exactly one button element with the provided text.
	// Supports <button>, <input type="button">, and <input type="submit">.
	FindByButton(text string) Selection

	// First finds the first element by CSS selector.
	First(selector string) Selection

	// FirstByXPath finds the first element by XPath selector.
	FirstByXPath(selector string) Selection

	// FirstByLink finds the first anchor element by its text content.
	FirstByLink(text string) Selection

	// FirstByLabel finds the first element by associated label text.
	FirstByLabel(text string) Selection

	// FirstByButton finds the first button element with the provided text.
	// Supports <button>, <input type="button">, and <input type="submit">.
	FirstByButton(text string) Selection

	// All finds zero or more elements by CSS selector.
	All(selector string) MultiSelection

	// AllByXPath finds zero or more elements by XPath selector.
	AllByXPath(selector string) MultiSelection

	// AllByLink finds zero or more anchor elements by their text content.
	AllByLink(text string) MultiSelection

	// AllByLabel finds zero or more elements by associated label text.
	AllByLabel(text string) MultiSelection

	// AllByButton finds zero or more button elements with the provided text.
	// Supports <button>, <input type="button">, and <input type="submit">.
	AllByButton(text string) MultiSelection
}
