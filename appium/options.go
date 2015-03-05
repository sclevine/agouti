package appium

import "github.com/sclevine/agouti"

func Desired(capabilities agouti.Capabilities) Option {
	return func(c *config) {
		c.desired = capabilities
	}
}

type config struct {
	desired agouti.Capabilities
}

type Option func(*config)

func (c config) merge(options []Option) *config {
	for _, option := range options {
		option(&c)
	}
	return &c
}
