package selection

import "fmt"

func (s *Selection) Text() (string, error) {
	element, err := s.getSelectedElement()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	text, err := element.GetText()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve text for '%s': %s", s, err)
	}
	return text, nil
}

func (s *Selection) Attribute(attribute string) (string, error) {
	element, err := s.getSelectedElement()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	value, err := element.GetAttribute(attribute)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve attribute value for '%s': %s", s, err)
	}
	return value, nil
}

func (s *Selection) CSS(property string) (string, error) {
	element, err := s.getSelectedElement()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	value, err := element.GetCSS(property)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve CSS property for '%s': %s", s, err)
	}
	return value, nil
}

func (s *Selection) Selected() (bool, error) {
	element, err := s.getSelectedElement()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	selected, err := element.IsSelected()
	if err != nil {
		return false, fmt.Errorf("failed to determine whether '%s' is selected: %s", s, err)
	}

	return selected, nil
}

func (s *Selection) Visible() (bool, error) {
	element, err := s.getSelectedElement()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	visible, err := element.IsDisplayed()
	if err != nil {
		return false, fmt.Errorf("failed to determine whether '%s' is visible: %s", s, err)
	}

	return visible, nil
}

func (s *Selection) Enabled() (bool, error) {
	element, err := s.getSelectedElement()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	enabled, err := element.IsEnabled()
	if err != nil {
		return false, fmt.Errorf("failed to determine whether '%s' is enabled: %s", s, err)
	}

	return enabled, nil
}
