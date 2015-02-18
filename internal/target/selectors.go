package target

import (
	"fmt"
	"strings"
)

type Selectors []Selector

func (s Selectors) AppendCSS(cssSelector string) Selectors {
	selector := Selector{Type: "css selector", Value: cssSelector}

	if s.canMergeCSS() {
		lastIndex := len(s) - 1
		selector.Value = s[lastIndex].Value + " " + selector.Value
		return appendSelector(s[:lastIndex], selector)
	}

	return appendSelector(s, selector)
}

func (s Selectors) canMergeCSS() bool {
	if len(s) == 0 {
		return false
	}
	last := s[len(s)-1]
	return last.Type == "css selector" && !last.Indexed && !last.Single
}

func (s Selectors) AppendXPath(xPathSelector string) Selectors {
	selector := Selector{Type: "xpath", Value: xPathSelector}
	return appendSelector(s, selector)
}

func (s Selectors) AppendLink(text string) Selectors {
	selector := Selector{Type: "link text", Value: text}
	return appendSelector(s, selector)
}

func (s Selectors) AppendLabeled(text string) Selectors {
	return s.AppendXPath(fmt.Sprintf(`//input[@id=(//label[normalize-space()="%s"]/@for)] | //label[normalize-space()="%s"]/input`, text, text))
}

func (s Selectors) AppendButton(text string) Selectors {
	return s.AppendXPath(fmt.Sprintf(`//input[@type="submit" or @type="button"][normalize-space(@value)="%s"] | //button[normalize-space()="%s"]`, text, text))
}

func (s Selectors) Single() Selectors {
	lastIndex := len(s) - 1
	if lastIndex < 0 {
		return nil
	}

	selector := s[lastIndex]
	selector.Single = true
	selector.Indexed = false
	return appendSelector(s[:lastIndex], selector)
}

func (s Selectors) At(index int) Selectors {
	lastIndex := len(s) - 1
	if lastIndex < 0 {
		return nil
	}

	selector := s[lastIndex]
	selector.Single = false
	selector.Indexed = true
	selector.Index = index
	return appendSelector(s[:lastIndex], selector)
}

func (s Selectors) String() string {
	var tags []string

	for _, selector := range s {
		tags = append(tags, selector.String())
	}

	return strings.Join(tags, " | ")
}

func appendSelector(selectors Selectors, selector Selector) Selectors {
	selectorsCopy := append(Selectors(nil), selectors...)
	return append(selectorsCopy, selector)
}
