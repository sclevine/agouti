package core_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/sclevine/agouti/core"
)

var _ = Describe("Cookie", func() {
	var cookie WebCookie

	BeforeEach(func() {
		cookie = Cookie("some-name", "some value")
	})

	It("should successfully encode all provided options into JSON", func() {
		cookie.Path("some/path").Domain("some.domain")
		cookie.Secure().HTTPOnly().Expiry(1000)
		json, err := cookie.JSON()
		Expect(err).NotTo(HaveOccurred())
		Expect(json).To(MatchJSON(`{
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
