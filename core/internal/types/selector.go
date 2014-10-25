package types

import "fmt"

type Selector struct {
	Using   string `json:"using"`
	Value   string `json:"value"`
	Index   int    `json:"-"`
	Indexed bool   `json:"-"`
	Single  bool   `json:"-"`
}

func (s Selector) String() string {
	var suffix string

	if s.Single {
		suffix = " [single]"
	} else if s.Indexed {
		suffix = fmt.Sprintf(" [%d]", s.Index)
	}

	switch s.Using {
	case "css selector":
		return fmt.Sprintf("CSS: %s%s", s.Value, suffix)
	case "xpath":
		return fmt.Sprintf("XPath: %s%s", s.Value, suffix)
	case "link text":
		return fmt.Sprintf(`Link: "%s"%s`, s.Value, suffix)
	default:
		return "Invalid selector"
	}
}
