package agouti

import "encoding/json"

// Cookie defines a cookie for use with *Page#SetCookie()
type Cookie map[string]interface{}

// NewCookie returns a Cookie instance with the provided name and value.
// All methods called on this instance will modify the original instance.
func NewCookie(name string, value string) Cookie {
	return Cookie{"name": name, "value": value}
}

// Path sets the cookie path - defaults to "/".
func (c Cookie) Path(path string) Cookie {
	c["path"] = path
	return c
}

// Domain sets the domain that the cookie is visible to.
func (c Cookie) Domain(domain string) Cookie {
	c["domain"] = domain
	return c
}

// SetSecure marks the cookie as a secure cookie
func (c Cookie) SetSecure() Cookie {
	c["secure"] = true
	return c
}

// SetHTTPOnly marks the cookie as HTTP-only
func (c Cookie) SetHTTPOnly() Cookie {
	c["httpOnly"] = true
	return c
}

// Expiry sets when the cookie expires in Unix time.
func (c Cookie) Expiry(expiry int64) Cookie {
	c["expiry"] = expiry
	return c
}

// JSON returns a JSON string representing the cookie.
func (c Cookie) JSON() (string, error) {
	cookieJSON, err := json.Marshal(c)
	return string(cookieJSON), err
}
