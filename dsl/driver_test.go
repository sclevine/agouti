package dsl_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/dsl"
)

var _ = Describe("Driver", func() {
	var (
		failMessage string
		failOffset  int
	)

	BeforeEach(func() {
		failMessage = ""
		RegisterAgoutiFailHandler(func(message string, callerSkip ...int) {
			failMessage = message
			failOffset = callerSkip[0]
			panic("Failed to catch test panic.")
		})
	})

	AfterEach(func() {
		RegisterAgoutiFailHandler(Fail)
	})

	Describe("Starting WebDrivers", func() {
		var originalPATH string

		BeforeEach(func() {
			originalPATH = os.Getenv("PATH")
		})

		AfterEach(func() {
			os.Setenv("PATH", originalPATH)
			StopWebDriver()
		})

		Describe("#StartPhantomJS", func() {
			It("should permit the creation and destruction of pages", func() {
				StartPhantomJS()
				page := CreatePage()
				Destroy(page)
			})

			It("should fail if a WebDriver is already started", func() {
				StartPhantomJS()
				Expect(StartPhantomJS).To(Panic())
				Expect(failMessage).To(Equal("WebDriver already started"))
				Expect(failOffset).To(Equal(2))
			})

			It("should fail if the WebDriver fails to start", func() {
				os.Setenv("PATH", "")
				Expect(StartPhantomJS).To(Panic())
				Expect(failMessage).To(Equal(`Agouti failure: failed to start service: failed to run command: exec: "phantomjs": executable file not found in $PATH`))
				Expect(failOffset).To(Equal(2))
			})
		})

		if os.Getenv("HEADLESS_ONLY") != "true" {
			Describe("#StartChromeDriver", func() {
				It("should permit the creation and destruction of pages", func() {
					StartChromeDriver()
					page := CreatePage()
					Destroy(page)
				})

				It("should fail if a WebDriver is already started", func() {
					StartSelenium()
					Expect(StartChromeDriver).To(Panic())
					Expect(failMessage).To(Equal("WebDriver already started"))
					Expect(failOffset).To(Equal(2))
				})

				It("should fail if the WebDriver fails to start", func() {
					os.Setenv("PATH", "")
					Expect(StartChromeDriver).To(Panic())
					Expect(failMessage).To(Equal(`Agouti failure: failed to start service: failed to run command: exec: "chromedriver": executable file not found in $PATH`))
					Expect(failOffset).To(Equal(2))
				})
			})

			Describe("#StartSelenium", func() {
				It("should permit the creation and destruction of pages", func() {
					StartSelenium()
					page := CreatePage("firefox")
					Destroy(page)
				})

				It("should fail if a WebDriver is already started", func() {
					StartChromeDriver()
					Expect(StartSelenium).To(Panic())
					Expect(failMessage).To(Equal("WebDriver already started"))
					Expect(failOffset).To(Equal(2))
				})

				It("should fail if the WebDriver fails to start", func() {
					os.Setenv("PATH", "")
					Expect(StartSelenium).To(Panic())
					Expect(failMessage).To(Equal(`Agouti failure: failed to start service: failed to run command: exec: "selenium-server": executable file not found in $PATH`))
					Expect(failOffset).To(Equal(2))
				})
			})
		}
	})

	Describe("#StopWebDriver", func() {
		It("should fail if the WebDriver is not already started", func() {
			Expect(StopWebDriver).To(Panic())
			Expect(failMessage).To(Equal("WebDriver not started"))
			Expect(failOffset).To(Equal(1))
		})
	})

	Describe("Pages", func() {
		Describe("#CreatePage", func() {
			It("should fail when there is no running WebDriver", func() {
				Expect(func() { CreatePage() }).To(Panic())
				Expect(failMessage).To(Equal("WebDriver not started"))
				Expect(failOffset).To(Equal(1))
			})
		})

		Describe("#CustomPage", func() {
			It("should fail when there is no running WebDriver", func() {
				Expect(func() { CustomPage(agouti.NewCapabilities()) }).To(Panic())
				Expect(failMessage).To(Equal("WebDriver not started"))
				Expect(failOffset).To(Equal(1))
			})
		})
	})
})
