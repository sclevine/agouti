package selection

import (
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

func (s *Selection) EqualsElement(other *Selection) (bool, error) {
	element, err := s.getSelectedElement()
	if err != nil {
		return false, fmt.Errorf("failed to select '%s': %s", s, err)
	}

	otherElement, err := other.getSelectedElement()
	if err != nil {
		return false, fmt.Errorf("failed to select '%s': %s", other, err)
	}

	equal, err := element.IsEqualTo(otherElement)
	if err != nil {
		return false, fmt.Errorf("failed to compare '%s' to '%s': %s", s, other, err)
	}

	return equal, nil
}
