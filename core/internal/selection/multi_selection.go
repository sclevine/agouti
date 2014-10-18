package selection

import (
	"fmt"
)

type MultiSelection struct {
	selection *Selection
}

func (m *MultiSelection) String() string {
	return m.selection.String() + " - All"
}

func (m *MultiSelection) Visible() (bool, error) {
	_, err := m.selection.getElements()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve elements with '%s': %s", m, err)
	}
	return false, nil

//	if len(elements) == 0 {
//		return false, fmt.Errorf("no elements ")
//	}
//
//	visible, err := element.IsDisplayed()
//	if err != nil {
//		return false, fmt.Errorf("failed to determine whether '%s' is visible: %s", s, err)
//	}
//
//	return visible, nil
}
