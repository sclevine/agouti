package page

import (
	"fmt"
	"github.com/sclevine/agouti/page/internal/webdriver"
	"strings"
)

type Selection interface {
	Find(selector string) Selection
	Selector() string
	Click() error
	Fill(text string) error
	Text() (string, error)
	Attribute(attribute string) (string, error)
	CSS(property string) (string, error)
	Check() error
	Selected() (bool, error)
	Select(text string) error
}

type selection struct {
	driver    driver
	selectors []string
}

func (s *selection) Find(selector string) Selection {
	return &selection{s.driver, append(s.selectors, selector)}
}

func (s *selection) Selector() string {
	return strings.Join(s.selectors, " ")
}

func (s *selection) Click() error {
	element, err := s.getSingleElement()
	if err != nil {
		return fmt.Errorf("failed to retrieve element with selector '%s': %s", s.Selector(), err)
	}

	if err := element.Click(); err != nil {
		return fmt.Errorf("failed to click on selector '%s': %s", s.Selector(), err)
	}
	return nil
}

func (s *selection) Fill(text string) error {
	element, err := s.getSingleElement()
	if err != nil {
		return fmt.Errorf("failed to retrieve element with selector '%s': %s", s.Selector(), err)
	}

	if err := element.Clear(); err != nil {
		return fmt.Errorf("failed to clear selector '%s': %s", s.Selector(), err)
	}

	if err := element.Value(text); err != nil {
		return fmt.Errorf("failed to enter text into selector '%s': %s", s.Selector(), err)
	}
	return nil
}

func (s *selection) Check() error {
	element, err := s.getSingleElement()
	if err != nil {
		return fmt.Errorf("failed to retrieve element with selector '%s': %s", s.Selector(), err)
	}

	elementType, err := element.GetAttribute("type")
	if err != nil {
		return fmt.Errorf("failed to retrieve type of selector '%s': %s", s.Selector(), err)
	}

	if elementType != "checkbox" {
		return fmt.Errorf("selector '%s' does not refer to a checkbox", s.Selector())
	}

	selected, err := element.Selected()
	if err != nil {
		return fmt.Errorf("failed to retrieve selected state of selector '%s': %s", s.Selector(), err)
	}

	if !selected {
		if err := element.Click(); err != nil {
			return fmt.Errorf("failed to check selector '%s': %s", s.Selector(), err)
		}
	}

	return nil
}

func (s *selection) Text() (string, error) {
	element, err := s.getSingleElement()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve element with selector '%s': %s", s.Selector(), err)
	}

	text, err := element.GetText()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve text for selector '%s': %s", s.Selector(), err)
	}
	return text, nil
}

func (s *selection) Attribute(attribute string) (string, error) {
	element, err := s.getSingleElement()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve element with selector '%s': %s", s.Selector(), err)
	}

	value, err := element.GetAttribute(attribute)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve attribute value for selector '%s': %s", s.Selector(), err)
	}
	return value, nil
}

func (s *selection) CSS(property string) (string, error) {
	element, err := s.getSingleElement()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve element with selector '%s': %s", s.Selector(), err)
	}

	value, err := element.GetCSS(property)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve CSS property for selector '%s': %s", s.Selector(), err)
	}
	return value, nil
}

func (s *selection) Selected() (bool, error) {
	element, err := s.getSingleElement()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve element with selector '%s': %s", s.Selector(), err)
	}

	selected, err := element.Selected()
	if err != nil {
		return false, fmt.Errorf("failed to determine whether selector '%s' is selected: %s", s.Selector(), err)
	}

	return selected, nil
}

func (s *selection) Select(text string) error {
	elements, err := s.driver.GetElements(s.Selector() + " option")
	if err != nil {
		return fmt.Errorf("failed to retrieve options for selector '%s': %s", s.Selector(), err)
	}

	for _, element := range elements {
		elementText, err := element.GetText()
		if err != nil {
			return fmt.Errorf("failed to retrieve option text for selector '%s': %s", s.Selector(), err)
		}

		if elementText == text {
			if err := element.Click(); err != nil {
				return fmt.Errorf(`failed to click on option with text "%s" for selector '%s': %s`, elementText, s.Selector(), err)
			}
			return nil
		}
	}

	return fmt.Errorf(`no options with text "%s" found for selector '%s'`, text, s.Selector())
}

func (s *selection) getSingleElement() (webdriver.Element, error) {
	elements, err := s.driver.GetElements(s.Selector())

	if err != nil {
		return nil, err
	}
	if len(elements) > 1 {
		return nil, fmt.Errorf("mutiple elements (%d) were selected", len(elements))
	}
	if len(elements) == 0 {
		return nil, fmt.Errorf("no element found")
	}

	return elements[0], nil
}
