package agouti

import "time"

type config struct {
	timeout time.Duration
	desired Capabilities
}

type Option func(*config)

// Timeout provides an option for specifying a timeout in seconds.
func Timeout(seconds int) Option {
	return func(c *config) {
		c.timeout = time.Duration(seconds) * time.Second
	}
}

// Desired provides an option for specifying desired WebDriver capabilities.
func Desired(capabilities Capabilities) Option {
	return func(c *config) {
		c.desired = capabilities
	}
}

func (c *config) apply(options []Option) *config {
	for _, option := range options {
		option(c)
	}
	return c
}

func getOptions(options []Option) *config {
	return (&config{}).apply(options)
}
