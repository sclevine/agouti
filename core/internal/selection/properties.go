package selection

import "fmt"

func (s *Selection) Text() (string, error) {
	element, err := s.Elements.GetExactlyOne(s.selectors)
	if err != nil {
		return "", fmt.Errorf("failed to select '%s': %s", s, err)
	}

	text, err := element.GetText()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve text for '%s': %s", s, err)
	}
	return text, nil
}

func (s *Selection) Active() (bool, error) {
	element, err := s.Elements.GetExactlyOne(s.selectors)
	if err != nil {
		return false, fmt.Errorf("failed to select '%s': %s", s, err)
	}

	activeElement, err := s.Session.GetActiveElement()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve active element: %s", err)
	}

	equal, err := element.IsEqualTo(activeElement)
	if err != nil {
		return false, fmt.Errorf("failed to compare selection to active element: %s", err)
	}

	return equal, nil
}

type propertyMethod func(element Element, property string) (string, error)

func (s *Selection) hasProperty(method propertyMethod, property, name string) (string, error) {
	element, err := s.Elements.GetExactlyOne(s.selectors)
	if err != nil {
		return "", fmt.Errorf("failed to select '%s': %s", s, err)
	}

	value, err := method(element, property)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve %s value for '%s': %s", name, s, err)
	}
	return value, nil
}

func (s *Selection) Attribute(attribute string) (string, error) {
	return s.hasProperty(Element.GetAttribute, attribute, "attribute")
}

func (s *Selection) CSS(property string) (string, error) {
	return s.hasProperty(Element.GetCSS, property, "CSS property")
}

type stateMethod func(element Element) (bool, error)

func (s *Selection) hasState(method stateMethod, name string) (bool, error) {
	elements, err := s.Elements.GetAtLeastOne(s.selectors)
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

func (s *Selection) Selected() (bool, error) {
	return s.hasState(Element.IsSelected, "selected")
}

func (s *Selection) Visible() (bool, error) {
	return s.hasState(Element.IsDisplayed, "visible")
}

func (s *Selection) Enabled() (bool, error) {
	return s.hasState(Element.IsEnabled, "enabled")
}
