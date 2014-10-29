package core

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
	// At finds the element at the provided index.
	At(index int) Selection

	// Single specifies that the selection must refer to exactly one element.
	//    selection.Find("#selector")
	// is equivalent to
	//    selection.All("#selector").Single()
	Single() Selection

	// All Selection methods are valid MultiSelection methods.
	Selection
}

type multiSelection struct {
	*selection
}

func (m *multiSelection) At(index int) Selection {
	last := len(m.selectors) - 1

	if last < 0 {
		return &selection{m.client, nil}
	}

	newSelector := m.selectors[last]
	newSelector.Indexed = true
	newSelector.Index = index

	return &selection{m.client, appendSelector(m.selectors[:last], newSelector)}
}

func (m *multiSelection) Single() Selection {
	last := len(m.selectors) - 1

	if last < 0 {
		return &selection{m.client, nil}
	}

	newSelector := m.selectors[last]
	newSelector.Single = true

	return &selection{m.client, appendSelector(m.selectors[:last], newSelector)}
}
