package selection

import (
	"fmt"

	"github.com/sclevine/agouti/core/internal/api"
)

type Selection struct {
	Client    apiClient
	Elements  elementRepository
	selectors []Selector
}

type apiClient interface {
	GetActiveElement() (*api.Element, error)
	DoubleClick() error
	MoveTo(element *api.Element, point api.Point) error
	Frame(frame *api.Element) error
}

type elementRepository interface {
	Get(selectors []Selector) ([]Element, error)
	GetAtLeastOne(selectors []Selector) ([]Element, error)
	GetExactlyOne(selectors []Selector) (Element, error)
}

func NewSelection(client *api.Client) *Selection {
	return &Selection{client, &ElementRepository{client}, nil}
}

func (s *Selection) AppendCSS(cssSelector string) *Selection {
	selector := Selector{Type: "css selector", Value: cssSelector}

	if s.canMergeCSS() {
		lastIndex := len(s.selectors) - 1
		selector.Value = s.selectors[lastIndex].Value + " " + selector.Value
		return &Selection{s.Client, s.Elements, appendSelector(s.selectors[:lastIndex], selector)}
	}

	return &Selection{s.Client, s.Elements, appendSelector(s.selectors, selector)}
}

func (s *Selection) canMergeCSS() bool {
	if len(s.selectors) == 0 {
		return false
	}
	last := s.selectors[len(s.selectors)-1]
	return last.Type == "css selector" && !last.Indexed && !last.Single
}

func (s *Selection) AppendXPath(xPathSelector string) *Selection {
	selector := Selector{Type: "xpath", Value: xPathSelector}
	return &Selection{s.Client, s.Elements, appendSelector(s.selectors, selector)}
}

func (s *Selection) AppendLink(text string) *Selection {
	selector := Selector{Type: "link text", Value: text}
	return &Selection{s.Client, s.Elements, appendSelector(s.selectors, selector)}
}

func (s *Selection) AppendLabeled(text string) *Selection {
	return s.AppendXPath(fmt.Sprintf(`//input[@id=(//label[normalize-space(text())="%s"]/@for)] | //label[normalize-space(text())="%s"]/input`, text, text))
}

func (s *Selection) Single() *Selection {
	lastIndex := len(s.selectors) - 1
	if lastIndex < 0 {
		return &Selection{s.Client, s.Elements, nil}
	}

	selector := s.selectors[lastIndex]
	selector.Single = true
	selector.Indexed = false
	return &Selection{s.Client, s.Elements, appendSelector(s.selectors[:lastIndex], selector)}
}

func (s *Selection) At(index int) *Selection {
	lastIndex := len(s.selectors) - 1
	if lastIndex < 0 {
		return &Selection{s.Client, s.Elements, nil}
	}

	selector := s.selectors[lastIndex]
	selector.Single = false
	selector.Indexed = true
	selector.Index = index
	return &Selection{s.Client, s.Elements, appendSelector(s.selectors[:lastIndex], selector)}
}

func appendSelector(selectors []Selector, selector Selector) []Selector {
	selectorsCopy := append([]Selector(nil), selectors...)
	return append(selectorsCopy, selector)
}
