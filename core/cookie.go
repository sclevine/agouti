package core

import "encoding/json"

// A Cookie instance defines a WebCookie for use with Page#SetCookie()
type WebCookie interface {
	Path(path string) WebCookie
	Domain(domain string) WebCookie
	Secure() WebCookie
	HTTPOnly() WebCookie
	Expiry(expiry int64) WebCookie

	// JSON returns a JSON string representing the cookie.
	JSON() (string, error)
}

// Cookie returns a WebCookie instance that can be passed to Page#SetCookie().
// All methods called on this instance will modify the original instance.
func Cookie(name string, value interface{}) WebCookie {
	return cookie{"name": name, "value": value}
}

type cookie map[string]interface{}

func (c cookie) Path(path string) WebCookie {
	c["path"] = path
	return c
}

func (c cookie) Domain(domain string) WebCookie {
	c["domain"] = domain
	return c
}

func (c cookie) Secure() WebCookie {
	c["secure"] = true
	return c
}

func (c cookie) HTTPOnly() WebCookie {
	c["httpOnly"] = true
	return c
}

func (c cookie) Expiry(expiry int64) WebCookie {
	c["expiry"] = expiry
	return c
}

func (c cookie) JSON() (string, error) {
	cookieJSON, err := json.Marshal(c)
	return string(cookieJSON), err
}
