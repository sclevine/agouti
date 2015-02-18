package core

import "github.com/sclevine/agouti"

// Cookie returns a agouti.Cookie instance with the provided name and value.
// This method, along with this entire package, is deprecated. It must return
// the new agouti.Cookie type to keep compatibility with the dsl package.
func Cookie(name string, value string) agouti.Cookie {
	return agouti.Cookie{"name": name, "value": value}
}
