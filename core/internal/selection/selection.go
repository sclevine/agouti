package selection

import (
	"fmt"

	"github.com/sclevine/agouti/core/internal/types"
)

type Selection struct {
	Client    Client
	selectors []types.Selector
}

type Client interface {
	DoubleClick() error
	MoveTo(element types.Element, point types.Point) error
	GetElement(selector types.Selector) (types.Element, error)
	GetElements(selector types.Selector) ([]types.Element, error)
}

func (s *Selection) Find(selector string) *Selection {
	return s.All(selector).Single()
}

func (s *Selection) FindByXPath(selector string) *Selection {
	return s.AllByXPath(selector).Single()
}

func (s *Selection) FindByLink(text string) *Selection {
	return s.AllByLink(text).Single()
}

func (s *Selection) FindByLabel(text string) *Selection {
	return s.AllByLabel(text).Single()
}

func (s *Selection) First(selector string) *Selection {
	return s.All(selector).At(0)
}

func (s *Selection) FirstByXPath(selector string) *Selection {
	return s.AllByXPath(selector).At(0)
}

func (s *Selection) FirstByLink(text string) *Selection {
	return s.AllByLink(text).At(0)
}

func (s *Selection) FirstByLabel(text string) *Selection {
	return s.AllByLabel(text).At(0)
}

func (s *Selection) All(selector string) *MultiSelection {
	last := len(s.selectors) - 1

	lastIsCSS := last >= 0 && s.selectors[last].Using == "css selector"
	if lastIsCSS && !s.selectors[last].Indexed && !s.selectors[last].Single {
		return s.mergedSelection(selector)
	}

	return s.subSelection("css selector", selector)
}

func (s *Selection) AllByXPath(selector string) *MultiSelection {
	return s.subSelection("xpath", selector)
}

func (s *Selection) AllByLink(text string) *MultiSelection {
	return s.subSelection("link text", text)
}

func (s *Selection) AllByLabel(text string) *MultiSelection {
	selector := fmt.Sprintf(`//input[@id=(//label[normalize-space(text())="%s"]/@for)] | //label[normalize-space(text())="%s"]/input`, text, text)
	return s.AllByXPath(selector)
}

func (s *Selection) subSelection(using, value string) *MultiSelection {
	newSelector := types.Selector{Using: using, Value: value}
	selection := &Selection{s.Client, appendSelector(s.selectors, newSelector)}
	return &MultiSelection{selection}
}

func (s *Selection) mergedSelection(value string) *MultiSelection {
	last := len(s.selectors) - 1
	newSelectorValue := s.selectors[last].Value + " " + value
	newSelector := types.Selector{Using: "css selector", Value: newSelectorValue}
	selection := &Selection{s.Client, appendSelector(s.selectors[:last], newSelector)}
	return &MultiSelection{selection}
}

func appendSelector(selectors []types.Selector, selector types.Selector) []types.Selector {
	selectorsCopy := append([]types.Selector(nil), selectors...)
	return append(selectorsCopy, selector)
}
