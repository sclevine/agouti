package selection

import (
	"fmt"

	"github.com/sclevine/agouti/api"
)

type Selection struct {
	Session interface {
		GetActiveElement() (*api.Element, error)
		DoubleClick() error
		MoveTo(element *api.Element, point api.Offset) error
		Frame(frame *api.Element) error
	}
	Elements interface {
		Get(selectors []Selector) ([]Element, error)
		GetAtLeastOne(selectors []Selector) ([]Element, error)
		GetExactlyOne(selectors []Selector) (Element, error)
	}
	selectors []Selector
}

func NewSelection(session *api.Session) *Selection {
	return &Selection{session, &ElementRepository{session}, nil}
}

func (s *Selection) AppendCSS(cssSelector string) *Selection {
	selector := Selector{Type: "css selector", Value: cssSelector}

	if s.canMergeCSS() {
		lastIndex := len(s.selectors) - 1
		selector.Value = s.selectors[lastIndex].Value + " " + selector.Value
		return &Selection{s.Session, s.Elements, appendSelector(s.selectors[:lastIndex], selector)}
	}

	return &Selection{s.Session, s.Elements, appendSelector(s.selectors, selector)}
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
	return &Selection{s.Session, s.Elements, appendSelector(s.selectors, selector)}
}

func (s *Selection) AppendLink(text string) *Selection {
	selector := Selector{Type: "link text", Value: text}
	return &Selection{s.Session, s.Elements, appendSelector(s.selectors, selector)}
}

func (s *Selection) AppendLabeled(text string) *Selection {
	return s.AppendXPath(fmt.Sprintf(`//input[@id=(//label[normalize-space()="%s"]/@for)] | //label[normalize-space()="%s"]/input`, text, text))
}

func (s *Selection) AppendButton(text string) *Selection {
	return s.AppendXPath(fmt.Sprintf(`//input[@type="submit" or @type="button"][normalize-space(@value)="%s"] | //button[normalize-space()="%s"]`, text, text))
}

func (s *Selection) Single() *Selection {
	lastIndex := len(s.selectors) - 1
	if lastIndex < 0 {
		return &Selection{s.Session, s.Elements, nil}
	}

	selector := s.selectors[lastIndex]
	selector.Single = true
	selector.Indexed = false
	return &Selection{s.Session, s.Elements, appendSelector(s.selectors[:lastIndex], selector)}
}

func (s *Selection) At(index int) *Selection {
	lastIndex := len(s.selectors) - 1
	if lastIndex < 0 {
		return &Selection{s.Session, s.Elements, nil}
	}

	selector := s.selectors[lastIndex]
	selector.Single = false
	selector.Indexed = true
	selector.Index = index
	return &Selection{s.Session, s.Elements, appendSelector(s.selectors[:lastIndex], selector)}
}

func appendSelector(selectors []Selector, selector Selector) []Selector {
	selectorsCopy := append([]Selector(nil), selectors...)
	return append(selectorsCopy, selector)
}
