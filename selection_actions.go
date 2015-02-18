package agouti

import (
	"fmt"

	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/element"
	"github.com/sclevine/agouti/internal/target"
)

type actionsFunc func(element.Element) error

func (s *Selection) forEachElement(actions actionsFunc) error {
	elements, err := s.elements.GetAtLeastOne(s.selectors)
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

// Click clicks on all of the elements that the selection refers to.
func (s *Selection) Click() error {
	return s.forEachElement(func(selectedElement element.Element) error {
		if err := selectedElement.Click(); err != nil {
			return fmt.Errorf("failed to click on '%s': %s", s, err)
		}
		return nil
	})
}

// DoubleClick double-clicks on all of the elements that the selection refers to.
func (s *Selection) DoubleClick() error {
	return s.forEachElement(func(selectedElement element.Element) error {
		if err := s.session.MoveTo(selectedElement.(*api.Element), nil); err != nil {
			return fmt.Errorf("failed to move mouse to '%s': %s", s, err)
		}
		if err := s.session.DoubleClick(); err != nil {
			return fmt.Errorf("failed to double-click on '%s': %s", s, err)
		}
		return nil
	})
}

// Fill fills all of the fields the selection refers to with the provided text.
func (s *Selection) Fill(text string) error {
	return s.forEachElement(func(selectedElement element.Element) error {
		if err := selectedElement.Clear(); err != nil {
			return fmt.Errorf("failed to clear '%s': %s", s, err)
		}
		if err := selectedElement.Value(text); err != nil {
			return fmt.Errorf("failed to enter text into '%s': %s", s, err)
		}
		return nil
	})
}

// Check checks all of the unchecked checkboxes that the selection refers to.
func (s *Selection) Check() error {
	return s.setChecked(true)
}

// Uncheck unchecks all of the checked checkboxes that the selection refers to.
func (s *Selection) Uncheck() error {
	return s.setChecked(false)
}

func (s *Selection) setChecked(checked bool) error {
	return s.forEachElement(func(selectedElement element.Element) error {
		elementType, err := selectedElement.GetAttribute("type")
		if err != nil {
			return fmt.Errorf("failed to retrieve type of '%s': %s", s, err)
		}

		if elementType != "checkbox" {
			return fmt.Errorf("'%s' does not refer to a checkbox", s)
		}

		elementChecked, err := selectedElement.IsSelected()
		if err != nil {
			return fmt.Errorf("failed to retrieve state of '%s': %s", s, err)
		}

		if elementChecked != checked {
			if err := selectedElement.Click(); err != nil {
				return fmt.Errorf("failed to click on '%s': %s", s, err)
			}
		}
		return nil
	})
}

// Select, when called on some number of <select> elements, will select all
// <option> elements under those <select> elements that match the provided text.
func (s *Selection) Select(text string) error {
	return s.forEachElement(func(selectedElement element.Element) error {
		optionXPath := fmt.Sprintf(`./option[normalize-space()="%s"]`, text)
		optionToSelect := target.Selector{Type: "xpath", Value: optionXPath}
		options, err := selectedElement.GetElements(optionToSelect.API())
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

// Submit submits all selected forms. The selection may refer to a form itself
// or any input element contained within a form.
func (s *Selection) Submit() error {
	return s.forEachElement(func(selectedElement element.Element) error {
		if err := selectedElement.Submit(); err != nil {
			return fmt.Errorf("failed to submit '%s': %s", s, err)
		}
		return nil
	})
}
