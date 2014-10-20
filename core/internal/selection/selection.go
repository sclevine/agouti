package selection

import (
	"errors"
	"fmt"
	"github.com/sclevine/agouti/core/internal/types"
	"strings"
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

func (s *Selection) At(index int) types.Selection {
	last := len(s.selectors) - 1

	if last < 0 {
		return &Selection{s.Client, nil}
	}

	old := s.selectors[last]
	newSelector := types.Selector{Using: old.Using, Value: old.Value, Index: index, Indexed: true}
	return &Selection{s.Client, appendSelector(s.selectors[:last], newSelector)}
}

func (s *Selection) Find(selector string) types.Selection {
	return s.All(selector).At(0)
}

func (s *Selection) FindByXPath(selector string) types.Selection {
	return s.AllByXPath(selector).At(0)
}

func (s *Selection) FindByLink(text string) types.Selection {
	return s.AllByLink(text).At(0)
}

func (s *Selection) FindByLabel(text string) types.Selection {
	return s.AllByLabel(text).At(0)
}

func (s *Selection) All(selector string) types.MultiSelection {
	last := len(s.selectors) - 1
	if last >= 0 && s.selectors[last].Using == "css selector" && !s.selectors[last].Indexed {
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

func (s *Selection) subSelection(using, value string) *Selection {
	newSelector := types.Selector{Using: using, Value: value}
	return &Selection{s.Client, appendSelector(s.selectors, newSelector)}
}

func (s *Selection) mergedSelection(value string) *Selection {
	last := len(s.selectors) - 1
	newSelectorValue := s.selectors[last].Value + " " + value
	newSelector := types.Selector{Using: "css selector", Value: newSelectorValue}
	return &Selection{s.Client, appendSelector(s.selectors[:last], newSelector)}
}

func appendSelector(selectors []types.Selector, selector types.Selector) []types.Selector {
	selectorsCopy := append([]types.Selector(nil), selectors...)
	return append(selectorsCopy, selector)
}

func (s *Selection) String() string {
	var tags []string

	for _, selector := range s.selectors {
		tags = append(tags, selector.String())
	}

	return strings.Join(tags, " | ")
}

func (s *Selection) Count() (int, error) {
	elements, err := s.getElements()
	if err != nil {
		return 0, fmt.Errorf("failed to select '%s': %s", s, err)
	}

	return len(elements), nil
}

func (s *Selection) EqualsElement(comparable interface{}) (bool, error) {
	element, err := s.getSelectedElement()
	if err != nil {
		return false, fmt.Errorf("failed to select '%s': %s", s, err)
	}

	selection, ok := comparable.(*Selection)
	if !ok {
		return false, errors.New("provided object is not a selection")
	}

	otherElement, err := selection.getSelectedElement()
	if err != nil {
		return false, fmt.Errorf("failed to select '%s': %s", comparable, err)
	}

	equal, err := element.IsEqualTo(otherElement)
	if err != nil {
		return false, fmt.Errorf("failed to compare '%s' to '%s': %s", s, comparable, err)
	}

	return equal, nil
}
