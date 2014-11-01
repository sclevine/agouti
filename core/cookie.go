package core

import "encoding/json"

// A WebCookie defines a cookie for use with Page#SetCookie()
type WebCookie interface {
	// Path sets the cookie path - defaults to "/".
	Path(path string) WebCookie

	// Domain sets the domain that the cookie is visible to.
	Domain(domain string) WebCookie

	// Secure marks the cookie as a secure cookie
	Secure() WebCookie

	// HTTPOnly marks the cookie as HTTP-only
	HTTPOnly() WebCookie

	// Expiry sets when the cookie expires in Unix time.
	Expiry(expiry int64) WebCookie

	// JSON returns a JSON string representing the cookie.
	JSON() (string, error)
}

// Cookie returns a WebCookie instance with the provided name and value.
// All methods called on this instance will modify the original instance.
func Cookie(name string, value string) WebCookie {
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
