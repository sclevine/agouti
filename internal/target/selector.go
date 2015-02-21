package target

import (
	"fmt"

	"github.com/sclevine/agouti/api"
)

const (
	CSS    = "CSS: %s"
	XPath  = "XPath: %s"
	Link   = `Link: "%s"`
	Label  = `Label: "%s"`
	Button = `Button: "%s"`

	labelXPath  = `//input[@id=(//label[normalize-space()="%s"]/@for)] | //label[normalize-space()="%[1]s"]/input`
	buttonXPath = `//input[@type="submit" or @type="button"][normalize-space(@value)="%s"] | //button[normalize-space()="%[1]s"]`
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

	return fmt.Sprintf(s.Type, s.Value) + suffix
}

func (s Selector) API() api.Selector {
	return api.Selector{Using: s.apiType(), Value: s.value()}
}
func (s Selector) apiType() string {
	switch s.Type {
	case XPath, Label, Button:
		return "xpath"
	case CSS:
		return "css selector"
	case Link:
		return "link text"
	}
	return "Invalid selector"
}

func (s Selector) value() string {
	switch s.Type {
	case Label:
		return fmt.Sprintf(labelXPath, s.Value)
	case Button:
		return fmt.Sprintf(buttonXPath, s.Value)
	}
	return s.Value
}
