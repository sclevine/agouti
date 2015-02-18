package agouti

import (
	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/element"
	"github.com/sclevine/agouti/internal/target"
)

type selectable struct {
	session   selectionSession
	selectors target.Selectors
}

type selectionSession interface {
	element.Client
	GetActiveElement() (*api.Element, error)
	DoubleClick() error
	MoveTo(element *api.Element, point api.Offset) error
	Frame(frame *api.Element) error
}

// Find finds exactly one element by CSS selector.
func (s *selectable) Find(selector string) *Selection {
	return newSelection(s.session, s.selectors.AppendCSS(selector).Single())
}

// FindByXPath finds exactly one element by XPath selector.
func (s *selectable) FindByXPath(selector string) *Selection {
	return newSelection(s.session, s.selectors.AppendXPath(selector).Single())
}

// FindByLink finds exactly one anchor element by its text content.
func (s *selectable) FindByLink(text string) *Selection {
	return newSelection(s.session, s.selectors.AppendLink(text).Single())
}

// FindByLabel finds exactly one element by associated label text.
func (s *selectable) FindByLabel(text string) *Selection {
	return newSelection(s.session, s.selectors.AppendLabeled(text).Single())
}

// FindByButton finds exactly one button element with the provided text.
// Supports <button>, <input type="button">, and <input type="submit">.
func (s *selectable) FindByButton(text string) *Selection {
	return newSelection(s.session, s.selectors.AppendButton(text).Single())
}

// First finds the first element by CSS selector.
func (s *selectable) First(selector string) *Selection {
	return newSelection(s.session, s.selectors.AppendCSS(selector).At(0))
}

// FirstByXPath finds the first element by XPath selector.
func (s *selectable) FirstByXPath(selector string) *Selection {
	return newSelection(s.session, s.selectors.AppendXPath(selector).At(0))
}

// FirstByLink finds the first anchor element by its text content.
func (s *selectable) FirstByLink(text string) *Selection {
	return newSelection(s.session, s.selectors.AppendLink(text).At(0))
}

// FirstByLabel finds the first element by associated label text.
func (s *selectable) FirstByLabel(text string) *Selection {
	return newSelection(s.session, s.selectors.AppendLabeled(text).At(0))
}

// FirstByButton finds the first button element with the provided text.
// Supports <button>, <input type="button">, and <input type="submit">.
func (s *selectable) FirstByButton(text string) *Selection {
	return newSelection(s.session, s.selectors.AppendButton(text).At(0))
}

// All finds zero or more elements by CSS selector.
func (s *selectable) All(selector string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.AppendCSS(selector))
}

// AllByXPath finds zero or more elements by XPath selector.
func (s *selectable) AllByXPath(selector string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.AppendXPath(selector))
}

// AllByLink finds zero or more anchor elements by their text content.
func (s *selectable) AllByLink(text string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.AppendLink(text))
}

// AllByLabel finds zero or more elements by associated label text.
func (s *selectable) AllByLabel(text string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.AppendLabeled(text))
}

// AllByButton finds zero or more button elements with the provided text.
// Supports <button>, <input type="button">, and <input type="submit">.
func (s *selectable) AllByButton(text string) *MultiSelection {
	return newMultiSelection(s.session, s.selectors.AppendButton(text))
}
