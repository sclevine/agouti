package page

import (
	"strings"
	"fmt"
)

type Selection interface {
	Within(selector string, bodies ...callable) Selection
	FinalSelection
}

type FinalSelection interface {
	ShouldContainText(text string)
	Selector() string
}

type selection struct {
	selectors []string
	page      *Page
}

func (s *selection) Within(selector string, bodies ...callable) Selection {
	subSelection := &selection{append(s.selectors, selector), s.page}
	for _, body := range bodies {
		body.Call(subSelection)
	}
	return subSelection
}

func (s *selection) Selector() string {
	return strings.Join(s.selectors, " ")
}

func (s *selection) ShouldContainText(text string) {
	// NOTE: return after failing in case Fail does not panic
	selector := s.Selector()
	elements, err := s.page.Driver.GetElements(selector)
	if err != nil {
		s.page.Fail("Failed to retrieve elements: " + err.Error(), 1)
		return
	}
	if len(elements) > 1 {
		s.page.Fail(fmt.Sprintf("Mutiple elements (%d) were selected.", len(elements)), 1)
		return
	}
	if len(elements) == 0 {
		s.page.Fail("No elements found.", 1)
		return
	}
	elementText, err := elements[0].GetText()
	if err != nil {
		s.page.Fail(fmt.Sprintf("Failed to retrieve text for selector '%s': %s", selector, err), 1)
		return
	}

	if !strings.Contains(elementText, text) {
		s.page.Fail(fmt.Sprintf("Failed to find text '%s' for selector '%s'.\nFound: '%s'", text, selector, elementText), 1)
		return
	}
}
