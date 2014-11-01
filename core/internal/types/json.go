package types

type JSON interface {
	JSON() (string, error)
}
