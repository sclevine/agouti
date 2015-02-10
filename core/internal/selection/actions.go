package selection

import (
	"fmt"

	"github.com/sclevine/agouti/api"
)

type actionsFunc func(Element) error

func (s *Selection) forEachElement(actions actionsFunc) error {
	elements, err := s.Elements.GetAtLeastOne(s.selectors)
	if err != nil {
		return fmt.Errorf("failed to select '%s': %s", s, err)
	}

	for _, element := range elements {
		if err := actions(element); err != nil {
			return err
		}
	}
	return nil
}

func (s *Selection) Click() error {
	return s.forEachElement(func(element Element) error {
		if err := element.Click(); err != nil {
			return fmt.Errorf("failed to click on '%s': %s", s, err)
		}
		return nil
	})
}

func (s *Selection) DoubleClick() error {
	return s.forEachElement(func(element Element) error {
		if err := s.Session.MoveTo(element.(*api.Element), nil); err != nil {
			return fmt.Errorf("failed to move mouse to '%s': %s", s, err)
		}
		if err := s.Session.DoubleClick(); err != nil {
			return fmt.Errorf("failed to double-click on '%s': %s", s, err)
		}
		return nil
	})
}

func (s *Selection) Fill(text string) error {
	return s.forEachElement(func(element Element) error {
		if err := element.Clear(); err != nil {
			return fmt.Errorf("failed to clear '%s': %s", s, err)
		}
		if err := element.Value(text); err != nil {
			return fmt.Errorf("failed to enter text into '%s': %s", s, err)
		}
		return nil
	})
}

func (s *Selection) Check() error {
	return s.setChecked(true)
}

func (s *Selection) Uncheck() error {
	return s.setChecked(false)
}

func (s *Selection) setChecked(checked bool) error {
	return s.forEachElement(func(element Element) error {
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
	})
}

func (s *Selection) Select(text string) error {
	return s.forEachElement(func(element Element) error {
		optionXPath := fmt.Sprintf(`./option[normalize-space()="%s"]`, text)
		optionToSelect := Selector{Type: "xpath", Value: optionXPath}
		options, err := element.GetElements(optionToSelect.API())
		if err != nil {
			return fmt.Errorf("failed to select specified option for some '%s': %s", s, err)
		}

		if len(options) == 0 {
			return fmt.Errorf(`no options with text "%s" found for some '%s'`, text, s)
		}

		for _, option := range options {
			if err := option.Click(); err != nil {
				return fmt.Errorf(`failed to click on option with text "%s" for some '%s': %s`, text, s, err)
			}
		}
		return nil
	})
}

func (s *Selection) Submit() error {
	return s.forEachElement(func(element Element) error {
		if err := element.Submit(); err != nil {
			return fmt.Errorf("failed to submit '%s': %s", s, err)
		}
		return nil
	})
}
