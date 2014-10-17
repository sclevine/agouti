package types

type Selection interface {
	Find(selector string) Selection
	FindXPath(selector string) Selection
	FindByLabel(text string) Selection
	At(index int) Selection
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
