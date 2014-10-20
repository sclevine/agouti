package types

import "fmt"

type Selector struct {
	Using   string `json:"using"`
	Value   string `json:"value"`
	Index   int    `json:"-"`
	Indexed bool   `json:"-"`
}

func (s Selector) String() string {
	var index string
	if s.Indexed {
		index = fmt.Sprintf(" [%d]", s.Index)
	}

	switch s.Using {
	case "css selector":
		return "CSS: " + s.Value + index
	case "xpath":
		return "XPath: " + s.Value + index
	case "link text":
		return fmt.Sprintf(`Link: "%s"`, s.Value) + index
	default:
		return "Invalid selector"
	}
}
