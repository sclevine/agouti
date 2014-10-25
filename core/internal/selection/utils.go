package selection

import (
	"errors"
	"fmt"
	"strings"
)

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

	var selection *Selection
	switch selectable := comparable.(type) {
	case *Selection:
		selection = selectable
	case *MultiSelection:
		selection = selectable.Selection
	default:
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
