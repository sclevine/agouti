package page

import (
	"github.com/onsi/ginkgo"
	"strings"
	"fmt"
)

type Selection interface {
	Within(selector string, bodies ...callable) Selection
	FinalSelection
}

type FinalSelection interface {
	ShouldContainText(text string)
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

func (s *selection) ShouldContainText(text string) {
	selector := strings.Join(s.selectors, " ")
	elements, err := s.page.Driver.GetElements(selector)
	if err != nil {
		ginkgo.Fail("Failed to retrieve elements.", 1)
	}
	if len(elements) > 1 {
		ginkgo.Fail("Mutiple items were selected.", 1)
	}
	if len(elements) == 0 {
		ginkgo.Fail("No items were selected.", 1)
	}
	elementText, err := elements[0].GetText()
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("Failed to retrieve text for selector %s.", selector), 1)
	}

	if !strings.Contains(elementText, text) {
		ginkgo.Fail(fmt.Sprintf("Failed to find text '%s' for selector '%s'.\nFound: %s ", text, selector, elementText), 1)
	}
}

// TODO: unit test both selection and page with fake driver and fake ginkgo.Fail
