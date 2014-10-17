package selection

import "fmt"

func (s *Selection) Click() error {
	element, err := s.getSelectedElement()
	if err != nil {
		return fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	if err := element.Click(); err != nil {
		return fmt.Errorf("failed to click on '%s': %s", s, err)
	}
	return nil
}

func (s *Selection) DoubleClick() error {
	element, err := s.getSelectedElement()
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
	element, err := s.getSelectedElement()
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
	element, err := s.getSelectedElement()
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
	element, err := s.getSelectedElement()
	if err != nil {
		return fmt.Errorf("failed to retrieve element with '%s': %s", s, err)
	}

	if err := element.Submit(); err != nil {
		return fmt.Errorf("failed to submit '%s': %s", s, err)
	}
	return nil
}
