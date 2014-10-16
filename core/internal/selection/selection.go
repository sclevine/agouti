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
}

type driver interface {
	GetElements(selector types.Selector) ([]types.Element, error)
	DoubleClick() error
	MoveTo(element types.Element, point types.Point) error
}

func (s *Selection) Find(selector string) types.Selection {
	last := len(s.selectors) - 1

	if last == -1 || s.selectors[last].Using != "css selector" {
		return &Selection{s.Driver, append(s.selectors, types.Selector{"css selector", selector})}
	}

	newSelector := types.Selector{"css selector", s.selectors[last].Value + " " + selector}
	newSelectors := append(append([]types.Selector(nil), s.selectors[:last]...), newSelector)
	return &Selection{s.Driver, newSelectors}
}

func (s *Selection) FindXPath(selector string) types.Selection {
	return &Selection{s.Driver, append(s.selectors, types.Selector{"xpath", selector})}
}

func (s *Selection) FindByLabel(text string) types.Selection {
	selector := fmt.Sprintf(`//input[@id=(//label[text()="%s"]/@for)] | //label[text()="%s"]/input`, text, text)
	return s.FindXPath(selector)
}

func (s *Selection) String() string {
	var tags []string

	for _, selector := range s.selectors {
		tags = append(tags, selector.String())
	}

	return strings.Join(tags, " | ")
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

func (s *Selection) getSingleElement() (types.Element, error) {
	elements, err := s.getElements()
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

func (s *Selection) Count() (int, error) {
	elements, err := s.getElements()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve elements for '%s': %s", s, err)
	}

	return len(elements), nil
}

func (s *Selection) Click() error {
	element, err := s.getSingleElement()
	if err != nil {
		return fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	if err := element.Click(); err != nil {
		return fmt.Errorf("failed to click on '%s': %s", s, err)
	}
	return nil
}

func (s *Selection) DoubleClick() error {
	element, err := s.getSingleElement()
	if err != nil {
		return fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	if err := s.Driver.MoveTo(element, nil); err != nil {
		return fmt.Errorf("failed to move mouse to '%s': %s", s, err)
	}

	if err := s.Driver.DoubleClick(); err != nil {
		return fmt.Errorf("failed to double-click on '%s': %s", s, err)
	}
	return nil
}

func (s *Selection) Fill(text string) error {
	element, err := s.getSingleElement()
	if err != nil {
		return fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	if err := element.Clear(); err != nil {
		return fmt.Errorf("failed to clear '%s': %s", s, err)
	}

	if err := element.Value(text); err != nil {
		return fmt.Errorf("failed to enter text into '%s': %s", s, err)
	}
	return nil
}

func (s *Selection) Check() error {
	return s.setChecked(true)
}

func (s *Selection) Uncheck() error {
	return s.setChecked(false)
}

func (s *Selection) setChecked(checked bool) error {
	element, err := s.getSingleElement()
	if err != nil {
		return fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	elementType, err := element.GetAttribute("type")
	if err != nil {
		return fmt.Errorf("failed to retrieve type of '%s': %s", s, err)
	}

	if elementType != "checkbox" {
		return fmt.Errorf("'%s' does not refer to a checkbox", s)
	}

	selected, err := element.IsSelected()
	if err != nil {
		return fmt.Errorf("failed to retrieve state of '%s': %s", s, err)
	}

	if selected != checked {
		if err := element.Click(); err != nil {
			return fmt.Errorf("failed to click on '%s': %s", s, err)
		}
	}

	return nil
}

func (s *Selection) Text() (string, error) {
	element, err := s.getSingleElement()
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
	element, err := s.getSingleElement()
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
	element, err := s.getSingleElement()
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
	element, err := s.getSingleElement()
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
	element, err := s.getSingleElement()
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
	element, err := s.getSingleElement()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	enabled, err := element.IsEnabled()
	if err != nil {
		return false, fmt.Errorf("failed to determine whether '%s' is enabled: %s", s, err)
	}

	return enabled, nil
}

func (s *Selection) Select(text string) error {
	elements, err := s.Find("option").(*Selection).getElements()
	if err != nil {
		return fmt.Errorf("failed to retrieve options for '%s': %s", s, err)
	}

	for _, element := range elements {
		elementText, err := element.GetText()
		if err != nil {
			return fmt.Errorf("failed to retrieve option text for '%s': %s", s, err)
		}

		if elementText == text {
			if err := element.Click(); err != nil {
				return fmt.Errorf(`failed to click on option with text "%s" for '%s': %s`, elementText, s, err)
			}
			return nil
		}
	}

	return fmt.Errorf(`no options with text "%s" found for '%s'`, text, s)
}

func (s *Selection) Submit() error {
	element, err := s.getSingleElement()
	if err != nil {
		return fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	if err := element.Submit(); err != nil {
		return fmt.Errorf("failed to submit '%s': %s", s, err)
	}
	return nil
}

func (s *Selection) EqualsElement(comparable interface{}) (bool, error) {
	element, err := s.getSingleElement()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	selection, ok := comparable.(*Selection)
	if !ok {
		return false, errors.New("provided object is not a selection")
	}

	otherElement, err := selection.getSingleElement()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve element with '%s': %s", comparable, err)
	}

	equal, err := element.IsEqualTo(otherElement)
	if err != nil {
		return false, fmt.Errorf("failed to compare '%s' to '%s': %s", s, comparable, err)
	}

	return equal, nil
}
