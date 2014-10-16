package selection

import (
	"fmt"
	"github.com/sclevine/agouti/core/internal/webdriver/types"
	"strings"
)

type Selection interface {
	Find(selector string) Selection
	FindXPath(selector string) Selection
	FindByLabel(text string) Selection
	String() string
	Count() (int, error)
	Click() error
	DoubleClick() error
	Fill(text string) error
	Text() (string, error)
	Attribute(attribute string) (string, error)
	CSS(property string) (string, error)
	Check() error
	Uncheck() error
	Selected() (bool, error)
	Visible() (bool, error)
	Enabled() (bool, error)
	Select(text string) error
	Submit() error
}

type selection struct {
	driver    driver
	selectors []types.Selector
}

type driver interface {
	GetElements(selector types.Selector) ([]types.Element, error)
	DoubleClick() error
	MoveTo(element types.Element, point types.Point) error
}

func New(driver driver, selector types.Selector) Selection {
	return &selection{driver, []types.Selector{selector}}
}

func (s *selection) Find(selector string) Selection {
	last := len(s.selectors) - 1

	if s.selectors[last].Using != "css selector" {
		return &selection{s.driver, append(s.selectors, types.Selector{"css selector", selector})}
	}

	newSelector := types.Selector{"css selector", s.selectors[last].Value + " " + selector}
	newSelectors := append(append([]types.Selector(nil), s.selectors[:last]...), newSelector)
	return &selection{s.driver, newSelectors}
}

func (s *selection) FindXPath(selector string) Selection {
	return &selection{s.driver, append(s.selectors, types.Selector{"xpath", selector})}
}

func (s *selection) FindByLabel(text string) Selection {
	selector := fmt.Sprintf(`//input[@id=(//label[text()="%s"]/@for)] | //label[text()="%s"]/input`, text, text)
	return s.FindXPath(selector)
}

func (s *selection) String() string {
	var tags []string

	for _, selector := range s.selectors {
		tags = append(tags, selector.String())
	}

	return strings.Join(tags, " | ")
}

func (s *selection) getElements() ([]types.Element, error) {
	lastElements, err := s.driver.GetElements(s.selectors[0])
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

func (s *selection) getSingleElement() (types.Element, error) {
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

func (s *selection) Count() (int, error) {
	elements, err := s.getElements()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve elements for '%s': %s", s, err)
	}

	return len(elements), nil
}

func (s *selection) Click() error {
	element, err := s.getSingleElement()
	if err != nil {
		return fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	if err := element.Click(); err != nil {
		return fmt.Errorf("failed to click on '%s': %s", s, err)
	}
	return nil
}

func (s *selection) DoubleClick() error {
	element, err := s.getSingleElement()
	if err != nil {
		return fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	if err := s.driver.MoveTo(element, nil); err != nil {
		return fmt.Errorf("failed to move mouse to '%s': %s", s, err)
	}

	if err := s.driver.DoubleClick(); err != nil {
		return fmt.Errorf("failed to double-click on '%s': %s", s, err)
	}
	return nil
}

func (s *selection) Fill(text string) error {
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

func (s *selection) Check() error {
	return s.setChecked(true)
}

func (s *selection) Uncheck() error {
	return s.setChecked(false)
}

func (s *selection) setChecked(checked bool) error {
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

func (s *selection) Text() (string, error) {
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

func (s *selection) Attribute(attribute string) (string, error) {
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

func (s *selection) CSS(property string) (string, error) {
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

func (s *selection) Selected() (bool, error) {
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

func (s *selection) Visible() (bool, error) {
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

func (s *selection) Enabled() (bool, error) {
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

func (s *selection) Select(text string) error {
	elements, err := s.Find("option").(*selection).getElements()
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

func (s *selection) Submit() error {
	element, err := s.getSingleElement()
	if err != nil {
		return fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	if err := element.Submit(); err != nil {
		return fmt.Errorf("failed to submit '%s': %s", s, err)
	}
	return nil
}
