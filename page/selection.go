package page

import (
	"fmt"
	"github.com/sclevine/agouti/webdriver"
	"strings"
)

type Selection interface {
	Should() FinalSelection
	ShouldNot() FinalSelection
	ShouldEventually() FinalSelection
	Within(selector string, bodies ...callable) Selection
	Click()
	Selector() string
}

type FinalSelection interface {
	ContainText(text string)
	Selector() string
}

type selection struct {
	selectors []string
	page      *page
	invert    bool
}

func (s *selection) Should() FinalSelection {
	return s
}

func (s *selection) ShouldNot() FinalSelection {
	s.invert = true
	return s
}

func (s *selection) ShouldEventually() FinalSelection {
	return &async{s}
}

func (s *selection) Within(selector string, bodies ...callable) Selection {
	subSelection := &selection{append(s.selectors, selector), s.page, false}
	for _, body := range bodies {
		body.Call(subSelection)
	}
	return subSelection
}

func (s *selection) Selector() string {
	return strings.Join(s.selectors, " ")
}

func (s *selection) ContainText(text string) {
	s.page.failer.Skip()
	element := s.getSingleElement()

	elementText, err := element.GetText()
	if err != nil {
		s.page.failer.Fail(fmt.Sprintf("Failed to retrieve text for selector '%s': %s", s.Selector(), err))
	}

	if strings.Contains(elementText, text) == s.invert {
		s.page.failer.Fail(fmt.Sprintf("%s text '%s' for selector '%s'.\nFound: '%s'", s.prefix(), text, s.Selector(), elementText))
	}
}

func (s *selection) Click() {
	s.page.failer.Skip()

	element := s.getSingleElement()

	if err := element.Click(); err != nil {
		s.page.failer.Fail(fmt.Sprintf("Failed to click on selector '%s': %s", s.Selector(), err))
	}
}

func (s *selection) getSingleElement() webdriver.Element {
	s.page.failer.Skip()

	selector := s.Selector()

	elements, err := s.page.driver.GetElements(selector)
	if err != nil {
		s.page.failer.Fail(fmt.Sprintf("Failed to retrieve element with selector '%s': %s", selector, err))
	}
	if len(elements) > 1 {
		s.page.failer.Fail(fmt.Sprintf("Mutiple elements (%d) with selector '%s' were selected.", len(elements), selector))
	}
	if len(elements) == 0 {
		s.page.failer.Fail(fmt.Sprintf("No element with selector '%s' found.", selector))
	}

	s.page.failer.UnSkip()
	return elements[0]
}

func (s *selection) prefix() string {
	if s.invert {
		return "Found"
	} else {
		return "Failed to find"
	}
}
