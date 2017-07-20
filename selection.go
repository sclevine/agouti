package agouti

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"path/filepath"

	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/crop"
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
	selectable
	elements elementRepository
	cropper  crop.Cropper
}

type elementRepository interface {
	Get() ([]element.Element, error)
	GetAtLeastOne() ([]element.Element, error)
	GetExactlyOne() (element.Element, error)
}

func newSelection(session apiSession, selectors target.Selectors) *Selection {
	return &Selection{
		selectable{session, selectors},
		&element.Repository{
			Client:    session,
			Selectors: selectors,
		},
		crop.CropperFunc(crop.Crop),
	}
}

// String returns a string representation of the selection, ex.
//    selection 'CSS: .some-class | XPath: //table [3] | Link "click me" [single]'
func (s *Selection) String() string {
	return fmt.Sprintf("selection '%s'", s.selectors)
}

// Elements returns a []*api.Element that can be used to send direct commands
// to WebDriver elements. See: https://code.google.com/p/selenium/wiki/JsonWireProtocol
func (s *Selection) Elements() ([]*api.Element, error) {
	elements, err := s.elements.Get()
	if err != nil {
		return nil, err
	}
	apiElements := []*api.Element{}
	for _, selectedElement := range elements {
		apiElements = append(apiElements, selectedElement.(*api.Element))
	}
	return apiElements, nil
}

// Count returns the number of elements that the selection refers to.
func (s *Selection) Count() (int, error) {
	elements, err := s.elements.Get()
	if err != nil {
		return 0, fmt.Errorf("failed to select elements from %s: %s", s, err)
	}

	return len(elements), nil
}

// EqualsElement returns whether or not two selections of exactly
// one element refer to the same element.
func (s *Selection) EqualsElement(other interface{}) (bool, error) {
	otherSelection, ok := other.(*Selection)
	if !ok {
		multiSelection, ok := other.(*MultiSelection)
		if !ok {
			return false, fmt.Errorf("must be *Selection or *MultiSelection")
		}
		otherSelection = &multiSelection.Selection
	}

	selectedElement, err := s.elements.GetExactlyOne()
	if err != nil {
		return false, fmt.Errorf("failed to select element from %s: %s", s, err)
	}

	otherElement, err := otherSelection.elements.GetExactlyOne()
	if err != nil {
		return false, fmt.Errorf("failed to select element from %s: %s", other, err)
	}

	equal, err := selectedElement.IsEqualTo(otherElement.(*api.Element))
	if err != nil {
		return false, fmt.Errorf("failed to compare %s to %s: %s", s, other, err)
	}

	return equal, nil
}

// MouseToElement moves the mouse over exactly one element in the selection.
func (s *Selection) MouseToElement() error {
	selectedElement, err := s.elements.GetExactlyOne()
	if err != nil {
		return fmt.Errorf("failed to select element from %s: %s", s, err)
	}

	if err := s.session.MoveTo(selectedElement.(*api.Element), nil); err != nil {
		return fmt.Errorf("failed to move mouse to element for %s: %s", s, err)
	}

	return nil
}

// Screenshot takes a screenshot of exactly one element
// and saves it to the provided filename.
// The provided filename may be an absolute or relative path.
func (s *Selection) Screenshot(filename string) error {
	selectedElement, err := s.elements.GetExactlyOne()
	if err != nil {
		return fmt.Errorf("failed to select element from %s: %s", s, err)
	}

	absFilePath, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to find absolute path for filename: %s", err)
	}

	screenshot, err := selectedElement.GetScreenshot()
	if err != nil {
		// Fallback to getting full size screenshot and cropping
		data, err := s.session.GetScreenshot()
		if err != nil {
			return fmt.Errorf("failed to retrieve screenshot: %s", err)
		}

		img, err := png.Decode(bytes.NewBuffer(data))
		if err != nil {
			return fmt.Errorf("failed to decode screenshot: %s", err)
		}

		x, y, width, height, err := selectedElement.GetRect()
		if err != nil {
			return fmt.Errorf("failed to retrieve bounds for selection: %s", err)
		}

		croppedImg, err := s.cropper.Crop(
			img, width, height,
			image.Point{X: x, Y: y})
		if err != nil {
			return fmt.Errorf("failed to crop screenshot: %s", err)
		}

		b := new(bytes.Buffer)
		err = png.Encode(b, croppedImg)
		if err != nil {
			return fmt.Errorf("failed to encode screenshot: %s", err)
		}

		screenshot = b.Bytes()
	}

	if err := ioutil.WriteFile(absFilePath, screenshot, 0666); err != nil {
		return fmt.Errorf("failed to save screenshot: %s", err)
	}

	return nil
}
