package agouti

import (
	"github.com/onsi/gomega"
	"github.com/onsi/ginkgo"
	"github.com/slevine/agouti/webdriver"
	"strings"
)

type Page struct {
	driver *webdriver.Driver
}

type Selection struct{
	selectors []string
	page *Page
}

func (p *Page) Within(selector string, bodies ...func(*Selection)) *Selection {
	selection := &Selection{[]string{selector}, p}
	for body := range bodies {
		body(selection)
	}
	return selection
}

func (s *Selection) Within(selector string, bodies ...func(*Selection)) *Selection {
	selection := &Selection{append(s.selectors, selector), s.page}
	for body := range bodies {
		body(selection)
	}
	return selection
}

func (s *Page) ShouldContainText(text string) {
	&Selection{[]string{"body"}, p}.ShouldContainText(text)
}

func (s *Selection) ShouldContainText(text string) {
	selector := strings.join(s.selectors, " ")
	elements := s.page.driver.GetElements(selector)
	if len(elements) > 1 {
		// TODO: how does one test a ginkgo fail
		ginkgo.Fail("Mutiple items were selected")
	}
	elementText := elements[0].GetText()
	gomega.Expect(elementText).To(gomega.ContainSubstring(text))
}
