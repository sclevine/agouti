package selection

import (
	"fmt"

	"github.com/sclevine/agouti/api"
)

type Selector struct {
	Type    string
	Value   string
	Index   int
	Indexed bool
	Single  bool
}

func (s Selector) String() string {
	var suffix string

	if s.Single {
		suffix = " [single]"
	} else if s.Indexed {
		suffix = fmt.Sprintf(" [%d]", s.Index)
	}

	switch s.Type {
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

func (s Selector) API() api.Selector {
	return api.Selector{Using: s.Type, Value: s.Value}
}
