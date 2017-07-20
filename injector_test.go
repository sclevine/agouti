package agouti

import (
	"github.com/sclevine/agouti/internal/crop"
	"github.com/sclevine/agouti/internal/target"
)

func NewTestSelection(session apiSession, elements elementRepository, firstSelector string, cropper crop.Cropper) *Selection {
	selector := target.Selector{Type: target.CSS, Value: firstSelector, Single: true}
	return &Selection{selectable{session, target.Selectors{selector}}, elements, cropper}
}

func NewTestMultiSelection(session apiSession, elements elementRepository, firstSelector string, cropper crop.Cropper) *MultiSelection {
	selector := target.Selector{Type: target.CSS, Value: firstSelector}
	selection := Selection{selectable{session, target.Selectors{selector}}, elements, cropper}
	return &MultiSelection{selection}
}

func NewTestPage(session apiSession) *Page {
	return &Page{selectable{session, nil}, nil}
}

func NewTestConfig() *config {
	return &config{}
}
