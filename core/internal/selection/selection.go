package selection

import (
	"fmt"

	"github.com/sclevine/agouti/core/internal/types"
)

type Selection struct {
	Client    types.Client
	selectors []types.Selector
}

func (s *Selection) AppendCSS(cssSelector string) *Selection {
	selector := types.Selector{Using: "css selector", Value: cssSelector}

	if s.canMergeCSS() {
		lastIndex := len(s.selectors) - 1
		selector.Value = s.selectors[lastIndex].Value + " " + selector.Value
		return &Selection{s.Client, appendSelector(s.selectors[:lastIndex], selector)}
	}

	return &Selection{s.Client, appendSelector(s.selectors, selector)}
}

func (s *Selection) canMergeCSS() bool {
	if len(s.selectors) == 0 {
		return false
	}
	last := s.selectors[len(s.selectors)-1]
	return last.Using == "css selector" && !last.Indexed && !last.Single
}

func (s *Selection) AppendXPath(xPathSelector string) *Selection {
	selector := types.Selector{Using: "xpath", Value: xPathSelector}
	return &Selection{s.Client, appendSelector(s.selectors, selector)}
}

func (s *Selection) AppendLink(text string) *Selection {
	selector := types.Selector{Using: "link text", Value: text}
	return &Selection{s.Client, appendSelector(s.selectors, selector)}
}

func (s *Selection) AppendLabeled(text string) *Selection {
	return s.AppendXPath(fmt.Sprintf(`//input[@id=(//label[normalize-space(text())="%s"]/@for)] | //label[normalize-space(text())="%s"]/input`, text, text))
}

func (s *Selection) Single() *Selection {
	lastIndex := len(s.selectors) - 1
	if lastIndex < 0 {
		return &Selection{s.Client, nil}
	}

	selector := s.selectors[lastIndex]
	selector.Single = true
	selector.Indexed = false
	return &Selection{s.Client, appendSelector(s.selectors[:lastIndex], selector)}
}

func (s *Selection) At(index int) *Selection {
	lastIndex := len(s.selectors) - 1
	if lastIndex < 0 {
		return &Selection{s.Client, nil}
	}

	selector := s.selectors[lastIndex]
	selector.Single = false
	selector.Indexed = true
	selector.Index = index
	return &Selection{s.Client, appendSelector(s.selectors[:lastIndex], selector)}
}

func appendSelector(selectors []types.Selector, selector types.Selector) []types.Selector {
	selectorsCopy := append([]types.Selector(nil), selectors...)
	return append(selectorsCopy, selector)
}
