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

	Describe("#Path", func() {
		It("should set a path", func() {
			json, _ := cookie.Path("/some/path").JSON()
			Expect(json).To(MatchJSON(`{"name": "some-name", "value": "some value", "path":"/some/path"}`))
		})
	})

	Describe("#Domain", func() {
		It("should set a domain", func() {
			json, _ := cookie.Domain("some.domain").JSON()
			Expect(json).To(MatchJSON(`{"name": "some-name", "value": "some value", "domain": "some.domain"}`))
		})
	})

	Describe("#Secure", func() {
		It("should mark the cookie as secure", func() {
			json, _ := cookie.Secure().JSON()
			Expect(json).To(MatchJSON(`{"name": "some-name", "value": "some value", "secure": true}`))
		})
	})

	Describe("#HTTPOnly", func() {
		It("should mark the cookie as HTTP-only", func() {
			json, _ := cookie.HTTPOnly().JSON()
			Expect(json).To(MatchJSON(`{"name": "some-name", "value": "some value", "httpOnly": true}`))
		})
	})

	Describe("#Expiry", func() {
		It("should set the cookie expiry", func() {
			json, _ := cookie.Expiry(1000).JSON()
			Expect(json).To(MatchJSON(`{"name": "some-name", "value": "some value", "expiry": 1000}`))
		})
	})

	Describe("#JSON", func() {
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

		Context("when the provided options cannot be converted to JSON", func() {
			It("should return an error", func() {
				cookie = Cookie("some-name", func() {})
				_, err := cookie.JSON()
				Expect(err).To(MatchError("json: unsupported type: func()"))
			})
		})
	})
})
