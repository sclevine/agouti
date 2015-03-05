package target

import (
	"fmt"

	"github.com/sclevine/agouti/api"
)

const (
	CSS        = "CSS: %s"
	XPath      = "XPath: %s"
	Link       = `Link: "%s"`
	Label      = `Label: "%s"`
	Button     = `Button: "%s"`
	A11yID     = "Accessibility ID: %s"
	AndroidAut = "Android UIAut.: %s"
	IOSAut     = "iOS UIAut.: %s"

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
	case CSS:
		return "css selector"
	case Link:
		return "link text"
	case A11yID:
		return "accessibility id"
	case AndroidAut:
		return "-android uiautomator"
	case IOSAut:
		return "-ios uiautomation"
	}
	return "xpath"
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
