package agouti

import "github.com/sclevine/agouti/internal/target"

func NewTestSelection(elements elementRepository, session selectionSession, firstSelector string) *Selection {
	selector := target.Selector{Type: "css selector", Value: firstSelector, Single: true}
	return &Selection{elements, selectable{session, target.Selectors{selector}}}
}

func NewTestMultiSelection(elements elementRepository, session selectionSession, firstSelector string) *MultiSelection {
	selector := target.Selector{Type: "css selector", Value: firstSelector}
	selection := Selection{elements, selectable{session, target.Selectors{selector}}}
	return &MultiSelection{selection}
}

func NewTestPage(session pageSession) *Page {
	return &Page{session, nil, selectable{}}
}
