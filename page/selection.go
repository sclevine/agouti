package page

import (
	"fmt"
	"strings"
	"github.com/sclevine/agouti/webdriver"
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
	page      *page
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
	element := s.getSingleElement()

	elementText, err := element.GetText()
	if err != nil {
		s.page.fail(fmt.Sprintf("Failed to retrieve text for selector '%s': %s", s.Selector(), err), 1)
	}

	if !strings.Contains(elementText, text) {
		s.page.fail(fmt.Sprintf("Failed to find text '%s' for selector '%s'.\nFound: '%s'", text, s.Selector(), elementText), 1)
	}
}

func (s *selection) getSingleElement() webdriver.Element {
	selector := s.Selector()

	elements, err := s.page.driver.GetElements(selector)
	if err != nil {
		s.page.fail("Failed to retrieve element: "+err.Error(), 2)
	}
	if len(elements) > 1 {
		s.page.fail(fmt.Sprintf("Mutiple elements (%d) were selected.", len(elements)), 2)
	}
	if len(elements) == 0 {
		s.page.fail("No element found.", 2)
	}
	return elements[0]
}
