package types

type Window interface {
	SetSize(height, width int) error
	SwitchTo() error
	Close() error
}
