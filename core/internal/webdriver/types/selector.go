package types

type Selector struct {
	Using string `json:"using"`
	Value string `json:"value"`
}

func (s Selector) String() string {
	switch s.Using {
	case "css selector":
		return "CSS: " + s.Value
	case "xpath":
		return "XPath: " + s.Value
	default:
		return "Invalid selector"
	}
}
