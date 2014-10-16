package types

type Element interface {
	GetID() string
	GetElements(selector Selector) ([]Element, error)
	GetText() (string, error)
	GetAttribute(attribute string) (string, error)
	GetCSS(property string) (string, error)
	IsSelected() (bool, error)
	IsDisplayed() (bool, error)
	IsEnabled() (bool, error)
	IsEqualTo(other Element) (bool, error)
	Click() error
	Clear() error
	Value(text string) error
	Submit() error
}
