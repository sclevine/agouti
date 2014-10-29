package core

import (
	"errors"
	"fmt"
	"strings"
)

func (s *selection) String() string {
	var tags []string

	for _, selector := range s.selectors {
		tags = append(tags, selector.String())
	}

	return strings.Join(tags, " | ")
}

func (s *selection) Count() (int, error) {
	elements, err := s.getElements()
	if err != nil {
		return 0, fmt.Errorf("failed to select '%s': %s", s, err)
	}

	return len(elements), nil
}

func (s *selection) EqualsElement(comparable interface{}) (bool, error) {
	element, err := s.getSelectedElement()
	if err != nil {
		return false, fmt.Errorf("failed to select '%s': %s", s, err)
	}

	var other *selection
	switch selectable := comparable.(type) {
	case *selection:
		other = selectable
	case *multiSelection:
		other = selectable.selection
	default:
		return false, errors.New("provided object is not a selection")
	}

	otherElement, err := other.getSelectedElement()
	if err != nil {
		return false, fmt.Errorf("failed to select '%s': %s", comparable, err)
	}

	equal, err := element.IsEqualTo(otherElement)
	if err != nil {
		return false, fmt.Errorf("failed to compare '%s' to '%s': %s", s, comparable, err)
	}

	return equal, nil
}
