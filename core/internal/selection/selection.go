package selection

import (
	"fmt"

	"github.com/sclevine/agouti/core/internal/types"
)

type Selection struct {
	Client    client
	selectors []types.Selector
}

type client interface {
	DoubleClick() error
	MoveTo(element types.Element, point types.Point) error
	retriever
}

func (s *Selection) Find(selector string) types.Selection {
	return s.All(selector).Single()
}

func (s *Selection) FindByXPath(selector string) types.Selection {
	return s.AllByXPath(selector).Single()
}

func (s *Selection) FindByLink(text string) types.Selection {
	return s.AllByLink(text).Single()
}

func (s *Selection) FindByLabel(text string) types.Selection {
	return s.AllByLabel(text).Single()
}

func (s *Selection) All(selector string) types.MultiSelection {
	last := len(s.selectors) - 1

	lastIsCSS := last >= 0 && s.selectors[last].Using == "css selector"
	if lastIsCSS && !s.selectors[last].Indexed && !s.selectors[last].Single {
		return s.mergedSelection(selector)
	}

	return s.subSelection("css selector", selector)
}

func (s *Selection) AllByXPath(selector string) types.MultiSelection {
	return s.subSelection("xpath", selector)
}

func (s *Selection) AllByLink(text string) types.MultiSelection {
	return s.subSelection("link text", text)
}

func (s *Selection) AllByLabel(text string) types.MultiSelection {
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
