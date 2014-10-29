package core

import (
	"fmt"

	"github.com/sclevine/agouti/core/internal/types"
)

func (s *selection) Text() (string, error) {
	element, err := s.getSelectedElement()
	if err != nil {
		return "", fmt.Errorf("failed to select '%s': %s", s, err)
	}

	text, err := element.GetText()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve text for '%s': %s", s, err)
	}
	return text, nil
}

type propertyMethod func(element types.Element, property string) (string, error)

func (s *selection) hasProperty(method propertyMethod, property, name string) (string, error) {
	element, err := s.getSelectedElement()
	if err != nil {
		return "", fmt.Errorf("failed to select '%s': %s", s, err)
	}

	value, err := method(element, property)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve %s value for '%s': %s", name, s, err)
	}
	return value, nil
}

func (s *selection) Attribute(attribute string) (string, error) {
	return s.hasProperty(types.Element.GetAttribute, attribute, "attribute")
}

func (s *selection) CSS(property string) (string, error) {
	return s.hasProperty(types.Element.GetCSS, property, "CSS property")
}

type stateMethod func(element types.Element) (bool, error)

func (s *selection) hasState(method stateMethod, name string) (bool, error) {
	elements, err := s.getSelectedElements()
	if err != nil {
		return false, fmt.Errorf("failed to select '%s': %s", s, err)
	}

	for _, element := range elements {
		pass, err := method(element)
		if err != nil {
			return false, fmt.Errorf("failed to determine whether some '%s' is %s: %s", s, name, err)
		}
		if !pass {
			return false, nil
		}
	}

	return true, nil
}

func (s *selection) Selected() (bool, error) {
	return s.hasState(types.Element.IsSelected, "selected")
}

func (s *selection) Visible() (bool, error) {
	return s.hasState(types.Element.IsDisplayed, "visible")
}

func (s *selection) Enabled() (bool, error) {
	return s.hasState(types.Element.IsEnabled, "enabled")
}
