package selection

import (
	"errors"
	"fmt"
	"github.com/sclevine/agouti/core/internal/types"
	"strings"
)

type Selection struct {
	Driver    driver
	selectors []types.Selector
	index     int
	indexed   bool
}

type driver interface {
	GetElements(selector types.Selector) ([]types.Element, error)
	DoubleClick() error
	MoveTo(element types.Element, point types.Point) error
}

func (s *Selection) getElements() ([]types.Element, error) {
	if len(s.selectors) == 0 {
		return nil, errors.New("empty selection")
	}

	lastElements, err := s.Driver.GetElements(s.selectors[0])
	if err != nil {
		return nil, err
	}

	for _, selector := range s.selectors[1:] {
		elements := []types.Element{}
		for _, element := range lastElements {
			subElements, err := element.GetElements(selector)
			if err != nil {
				return nil, err
			}
			elements = append(elements, subElements...)
		}
		lastElements = elements
	}
	return lastElements, nil
}

func (s *Selection) getSelectedElement() (types.Element, error) {
	elements, err := s.getElements()
	if err != nil {
		return nil, err
	}

	if len(elements) == 0 {
		return nil, fmt.Errorf("no element found")
	}

	if !s.indexed && len(elements) > 1  {
		return nil, fmt.Errorf("mutiple elements (%d) were selected", len(elements))
	}

	if s.index >= len(elements) {
		return nil, fmt.Errorf("element index (%d) out of range (>%d)", s.index, len(elements)-1)
	}

	return elements[s.index], nil
}

// TODO: fix mid-selection index
func (s *Selection) At(index int) types.Selection {
	return &Selection{s.Driver, s.selectors, index, true}
}

func (s *Selection) Find(selector string) types.Selection {
	last := len(s.selectors) - 1

	// TODO: fix double append bug!
	if last == -1 || s.selectors[last].Using != "css selector" {
		newSelector := types.Selector{Using: "css selector", Value: selector}
		return &Selection{s.Driver, append(s.selectors, newSelector), s.index, s.indexed}
	}

	newSelectorValue := s.selectors[last].Value + " " + selector
	newSelector := types.Selector{Using: "css selector", Value: newSelectorValue}
	newSelectors := append(append([]types.Selector(nil), s.selectors[:last]...), newSelector)
	return &Selection{s.Driver, newSelectors, s.index, s.indexed}
}

func (s *Selection) FindXPath(selector string) types.Selection {
	newSelector := types.Selector{Using: "xpath", Value: selector}
	return &Selection{s.Driver, append(s.selectors, newSelector), s.index, s.indexed}
}

func (s *Selection) FindByLabel(text string) types.Selection {
	selector := fmt.Sprintf(`//input[@id=(//label[normalize-space(text())="%s"]/@for)] | //label[normalize-space(text())="%s"]/input`, text, text)
	return s.FindXPath(selector)
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
		return 0, fmt.Errorf("failed to retrieve elements for '%s': %s", s, err)
	}

	return len(elements), nil
}

func (s *Selection) EqualsElement(comparable interface{}) (bool, error) {
	element, err := s.getSelectedElement()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	selection, ok := comparable.(*Selection)
	if !ok {
		return false, errors.New("provided object is not a selection")
	}

	otherElement, err := selection.getSelectedElement()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve element with '%s': %s", comparable, err)
	}

	equal, err := element.IsEqualTo(otherElement)
	if err != nil {
		return false, fmt.Errorf("failed to compare '%s' to '%s': %s", s, comparable, err)
	}

	return equal, nil
}
