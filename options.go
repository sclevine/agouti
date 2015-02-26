package agouti

import "time"

type config struct {
	timeout          time.Duration
	desired          Capabilities
	browser          string
	rejectInvalidSSL bool
}

// An Option specifies configuration for a new WebDriver or Page.
type Option func(*config)

// Browser provides an Option for specifying a browser.
func Browser(name string) Option {
	return func(c *config) {
		c.browser = name
	}
}

// Timeout provides an Option for specifying a timeout in seconds.
func Timeout(seconds int) Option {
	return func(c *config) {
		c.timeout = time.Duration(seconds) * time.Second
	}
}

// Desired provides an Option for specifying desired WebDriver Capabilities.
func Desired(capabilities Capabilities) Option {
	return func(c *config) {
		c.desired = capabilities
	}
}

// RejectInvalidSSL is an Option specifying that the WebDriver should reject
// invalid SSL certificates. All WebDrivers should accept invalid SSL certificates
// by default. See: http://www.w3.org/TR/webdriver/#invalid-ssl-certificates
func RejectInvalidSSL(c *config) {
	c.rejectInvalidSSL = true
}

func (c config) merge(options []Option) *config {
	for _, option := range options {
		option(&c)
	}
	return &c
}

func (c *config) capabilities() Capabilities {
	merged := Capabilities{"acceptSslCerts": true}
	for feature, value := range c.desired {
		merged[feature] = value
	}
	if c.browser != "" {
		merged.Browser(c.browser)
	}
	if c.rejectInvalidSSL {
		merged.Without("acceptSslCerts")
	}
	return merged
}
