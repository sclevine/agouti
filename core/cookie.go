package core

import "net/http"

// Cookie returns a *http.Cookie instance with the provided name and value.
// This method, along with this entire package, is deprecated. It must return
// an http.Cookie type to keep compatibility with the dsl package.
func Cookie(name string, value string) *http.Cookie {
	return &http.Cookie{Name: name, Value: value}
}
