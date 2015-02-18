package agouti_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti"
)

var _ = Describe("Cookie", func() {
	var cookie Cookie

	BeforeEach(func() {
		cookie = NewCookie("some-name", "some value")
	})

	It("should successfully encode all provided options into JSON", func() {
		cookie.Path("some/path").Domain("some.domain")
		cookie.SetSecure().SetHTTPOnly().Expiry(1000)
		Expect(cookie.JSON()).To(MatchJSON(`{
			"name": "some-name",
			"value": "some value",
			"path": "some/path",
			"domain": "some.domain",
			"secure": true,
			"httpOnly": true,
			"expiry": 1000
		}`))
	})
})
