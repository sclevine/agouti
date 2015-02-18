package agouti

import (
	"fmt"
	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/element"
	"github.com/sclevine/agouti/internal/target"
)

// Selection instances refer to a selection of elements.
// All Selection methods are also MultiSelection methods.
//
// Methods that take selectors apply their selectors to each element in the
// selection they are called on. If the selection they are called on refers to multiple
// elements, the resulting selection will refer to at least that many elements.
//
// Examples:
//
//    selection.Find("table").All("tr").At(2).First("td input[type=checkbox]").Check()
// Checks the first checkbox in the third row of the only table.
//    selection.Find("table").All("tr").Find("td").All("input[type=checkbox]").Check()
// Checks all checkboxes in the first-and-only cell of each row in the only table.
type Selection struct {
	elements elementRepository
	selectable
}

type elementRepository interface {
	Get(selectors target.Selectors) ([]element.Element, error)
	GetAtLeastOne(selectors target.Selectors) ([]element.Element, error)
	GetExactlyOne(selectors target.Selectors) (element.Element, error)
}

func newSelection(session selectionSession, selectors target.Selectors) *Selection {
	return &Selection{&element.Repository{session}, selectable{session, selectors}}
}

// String returns a string representation of the selection, ex.
//    CSS: .some-class | XPath: //table [3] | Link "click me" [single]
func (s *Selection) String() string {
	return s.selectors.String()
}

// Count returns the number of elements that the selection refers to.
func (s *Selection) Count() (int, error) {
	elements, err := s.elements.Get(s.selectors)
	if err != nil {
		return 0, fmt.Errorf("failed to select '%s': %s", s, err)
	}

	return len(elements), nil
}

// EqualsElement returns whether or not two selections of exactly
// one element each refer to the same element.
func (s *Selection) EqualsElement(other interface{}) (bool, error) {
	var otherSelection *Selection

	otherSelection, ok := other.(*Selection)
	if !ok {
		multiSelection, ok := other.(*MultiSelection)
		if !ok {
			return false, fmt.Errorf("must be *Selection or *MultiSelection")
		}
		otherSelection = &multiSelection.Selection
	}

	selectedElement, err := s.elements.GetExactlyOne(s.selectors)
	if err != nil {
		return false, fmt.Errorf("failed to select '%s': %s", s, err)
	}

	otherElement, err := otherSelection.elements.GetExactlyOne(s.selectors)
	if err != nil {
		return false, fmt.Errorf("failed to select '%s': %s", other, err)
	}

	equal, err := selectedElement.IsEqualTo(otherElement.(*api.Element))
	if err != nil {
		return false, fmt.Errorf("failed to compare '%s' to '%s': %s", s, other, err)
	}

	return equal, nil
}
