package page

import (
	"fmt"
	"github.com/sclevine/agouti/webdriver"
	"strings"
	"time"
)

type Selection interface {
	Should() FinalSelection
	ShouldNot() FinalSelection
	ShouldEventually(timing ...time.Duration) FinalSelection
	Within(selector string, bodies ...callable) Selection
	Click()
	Selector() string
}

type FinalSelection interface {
	ContainText(text string)
	Selector() string
}

type selection struct {
	driver    driver
	failer    failer
	selectors []string
	invert    bool
}

func (s *selection) Should() FinalSelection {
	return s
}

func (s *selection) ShouldNot() FinalSelection {
	s.invert = true
	return s
}

func (s *selection) ShouldEventually(timing ...time.Duration) FinalSelection {
	if len(timing) > 1 {
		return &async{s, timing[0], timing[1]}
	} else if len(timing) == 1 {
		return &async{s, timing[0], 100 * time.Millisecond}
	}

	return &async{s, 2 * time.Second, 100 * time.Millisecond}
}

func (s *selection) Within(selector string, bodies ...callable) Selection {
	subSelection := &selection{s.driver, s.failer, append(s.selectors, selector), false}
	for _, body := range bodies {
		body.Call(subSelection)
	}
	return subSelection
}

func (s *selection) Selector() string {
	return strings.Join(s.selectors, " ")
}

func (s *selection) ContainText(text string) {
	s.failer.Down()
	element := s.getSingleElement()

	elementText, err := element.GetText()
	if err != nil {
		s.failer.Fail(fmt.Sprintf("Failed to retrieve text for selector '%s': %s", s.Selector(), err))
	}

	if strings.Contains(elementText, text) == s.invert {
		s.failer.Fail(fmt.Sprintf("%s text '%s' for selector '%s'.\nFound: '%s'", s.prefix(), text, s.Selector(), elementText))
	}
	s.failer.Up()
}

func (s *selection) Click() {
	s.failer.Down()
	element := s.getSingleElement()

	if err := element.Click(); err != nil {
		s.failer.Fail(fmt.Sprintf("Failed to click on selector '%s': %s", s.Selector(), err))
	}
	s.failer.Up()
}

func (s *selection) getSingleElement() webdriver.Element {
	s.failer.Down()
	selector := s.Selector()

	elements, err := s.driver.GetElements(selector)
	if err != nil {
		s.failer.Fail(fmt.Sprintf("Failed to retrieve element with selector '%s': %s", selector, err))
	}
	if len(elements) > 1 {
		s.failer.Fail(fmt.Sprintf("Mutiple elements (%d) with selector '%s' were selected.", len(elements), selector))
	}
	if len(elements) == 0 {
		s.failer.Fail(fmt.Sprintf("No element with selector '%s' found.", selector))
	}
	s.failer.Up()

	return elements[0]
}

func (s *selection) prefix() string {
	if s.invert {
		return "Found"
	} else {
		return "Failed to find"
	}
}
