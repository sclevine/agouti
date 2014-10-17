package types

import "fmt"

type Selector struct {
	Using string `json:"using"`
	Value string `json:"value"`
	Index int `json:"-"`
	Indexed bool `json:"-"`
}

func (s Selector) String() string {
	text := s.Value
	if s.Indexed {
		text = text + fmt.Sprintf(" [%d]", s.Index)
	}

	switch s.Using {
	case "css selector":
		return "CSS: " + text
	case "xpath":
		return "XPath: " + text
	default:
		return "Invalid selector"
	}
}
