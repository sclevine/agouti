package types

type Selection interface {
	Find(selector string) Selection
	FindByXPath(selector string) Selection
	FindByLink(text string) Selection
	FindByLabel(text string) Selection
	All(selector string) MultiSelection
	AllByXPath(selector string) MultiSelection
	AllByLink(text string) MultiSelection
	AllByLabel(text string) MultiSelection
	String() string
	Count() (int, error)
	Click() error
	DoubleClick() error
	Fill(text string) error
	Text() (string, error)
	Attribute(attribute string) (string, error)
	CSS(property string) (string, error)
	Check() error
	Uncheck() error
	Selected() (bool, error)
	Visible() (bool, error)
	Enabled() (bool, error)
	Select(text string) error
	Submit() error
	EqualsElement(comparable interface{}) (bool, error)
}

type MultiSelection interface {
	Selection
	At(index int) Selection
}
