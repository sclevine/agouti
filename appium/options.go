package appium

import "github.com/sclevine/agouti"

func Desired(capabilities agouti.Capabilities) Option {
	return func(c *config) {
		c.desired = capabilities
	}
}

// Debug is used to configure a WebDriver in debug mode.
func Debug(state bool) Option {
	return func(c *config) {
		c.debug = true
	}
}

type config struct {
	desired agouti.Capabilities
	debug   bool
}

type Option func(*config)

func (c config) merge(options []Option) *config {
	for _, option := range options {
		option(&c)
	}
	return &c
}
