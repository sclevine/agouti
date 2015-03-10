package integration_test

import (
	"image/png"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("integration tests", func() {
	itShouldBehaveLikeAPage("PhantomJS", phantomDriver.NewPage)
	if !headlessOnly {
		itShouldBehaveLikeAPage("ChromeDriver", chromeDriver.NewPage)
		itShouldBehaveLikeAPage("Firefox", seleniumDriver.NewPage)
	}
})

type pageFunc func(...agouti.Option) (*agouti.Page, error)

func itShouldBehaveLikeAPage(name string, newPage pageFunc) {
	Describe("integration test for "+name, func() {
		var (
			page      *agouti.Page
			server    *httptest.Server
			submitted bool
		)

		BeforeEach(func() {
			server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				if request.Method == "POST" {
					submitted = true
				}
				html, _ := ioutil.ReadFile("test_page.html")
				response.Write(html)
			}))

			var err error
			page, err = newPage()
			Expect(err).NotTo(HaveOccurred())

			Expect(page.Size(640, 480)).To(Succeed())
			Expect(page.Navigate(server.URL)).To(Succeed())
		})

		AfterEach(func() {
			Expect(page.Destroy()).To(Succeed())
			server.Close()
		})

		Describe("Selection interactions", func() {
			It("should support asserting on element identity", func() {
				By("asserting on an element's existence", func() {
					Expect(page.Find("header")).To(BeFound())
					Expect(page.Find("header")).To(HaveCount(1))
					Expect(page.Find("not-a-header")).NotTo(BeFound())
				})

				By("comparing two selections for equality", func() {
					Expect(page.Find("#some_element")).To(EqualElement(page.FindByXPath("//div[@class='some-element']")))
				})
			})

			It("should support selecting elements", func() {
				By("finding an element by selection index", func() {
					Expect(page.All("option").At(0)).To(HaveText("first option"))
					Expect(page.All("select").At(1).First("option")).To(HaveText("third option"))
				})

				By("finding an element by chained selectors", func() {
					Expect(page.Find("header").Find("h1")).To(HaveText("Title"))
					Expect(page.Find("header").FindByXPath("//h1")).To(HaveText("Title"))
				})

				By("finding an element by link text", func() {
					Expect(page.FindByLink("Click Me").Attribute("href")).To(HaveSuffix("#new_page"))
				})

				By("finding an element by label text", func() {
					Expect(page.FindByLabel("Some Label")).To(HaveAttribute("value", "some labeled value"))
					Expect(page.FindByLabel("Some Container Label")).To(HaveAttribute("value", "some embedded value"))
				})

				By("finding an element by button text", func() {
					Expect(page.FindByButton("Some Button")).To(HaveAttribute("name", "some button name"))
					Expect(page.FindByButton("Some Input Button")).To(HaveAttribute("type", "button"))
					Expect(page.FindByButton("Some Submit Button")).To(HaveAttribute("type", "submit"))
				})

				By("finding an element by class", func() {
					Expect(page.FindByClass("some-element")).To(HaveAttribute("id", "some_element"))
				})

				By("finding an element by ID", func() {
					Expect(page.FindByID("some_element")).To(HaveAttribute("class", "some-element"))
				})

				By("finding multiple elements", func() {
					Expect(page.All("select").All("option")).To(BeVisible())
					Expect(page.All("h1,h2")).NotTo(BeVisible())
				})
			})

			It("should support retrieving element properties", func() {
				By("asserting on element text", func() {
					Expect(page.Find("header")).To(HaveText("Title"))
					Expect(page.Find("header")).NotTo(HaveText("Not-Title"))
					Expect(page.Find("header")).To(MatchText("T.+e"))
					Expect(page.Find("header")).NotTo(MatchText("X.+e"))
				})

				By("asserting on whether elements are active", func() {
					Expect(page.Find("#labeled_field")).NotTo(BeActive())
					Expect(page.Find("#labeled_field").Click()).To(Succeed())
					Expect(page.Find("#labeled_field")).To(BeActive())
				})

				By("asserting on element attributes", func() {
					Expect(page.Find("#some_checkbox")).To(HaveAttribute("type", "checkbox"))
				})

				By("asserting on element CSS", func() {
					Expect(page.Find("#some_element")).To(HaveCSS("color", "rgba(0, 0, 255, 1)"))
					Expect(page.Find("#some_element")).To(HaveCSS("color", "rgb(0, 0, 255)"))
					Expect(page.Find("#some_element")).To(HaveCSS("color", "blue"))
				})

				By("asserting on whether elements are selected", func() {
					Expect(page.Find("#some_checkbox")).NotTo(BeSelected())
					Expect(page.Find("#some_selected_checkbox")).To(BeSelected())
				})

				By("asserting on element visibility", func() {
					Expect(page.Find("header h1")).To(BeVisible())
					Expect(page.Find("header h2")).NotTo(BeVisible())
				})

				By("asserting on whether elements are enabled", func() {
					Expect(page.Find("#some_checkbox")).To(BeEnabled())
					Expect(page.Find("#some_disabled_checkbox")).NotTo(BeEnabled())
				})
			})

			It("should support element actions", func() {
				By("clicking on an element", func() {
					checkbox := page.Find("#some_checkbox")
					Expect(checkbox.Click()).To(Succeed())
					Expect(checkbox).To(BeSelected())
					Expect(checkbox.Click()).To(Succeed())
					Expect(checkbox).NotTo(BeSelected())
				})

				By("double-clicking on an element", func() {
					selection := page.Find("#double_click")
					Expect(selection.DoubleClick()).To(Succeed())
					Expect(selection).To(HaveText("double-click success"))
				})

				By("filling out an element", func() {
					Expect(page.Find("#some_input").Fill("some other value")).To(Succeed())
					Expect(page.Find("#some_input")).To(HaveAttribute("value", "some other value"))
				})

				// NOTE: PhantomJS regression causes crash on file upload
				if name != "PhantomJS" {
					By("uploading a file", func() {
						Expect(page.Find("#file_picker").UploadFile("test_page.html")).To(Succeed())
						var result string
						Expect(page.RunScript("return document.getElementById('file_picker').value;", nil, &result)).To(Succeed())
						Expect(result).To(HaveSuffix("test_page.html"))
					})
				}

				By("checking and unchecking a checkbox", func() {
					checkbox := page.Find("#some_checkbox")
					Expect(checkbox.Uncheck()).To(Succeed())
					Expect(checkbox).NotTo(BeSelected())
					Expect(checkbox.Check()).To(Succeed())
					Expect(checkbox).To(BeSelected())
					Expect(checkbox.Uncheck()).To(Succeed())
					Expect(checkbox).NotTo(BeSelected())
				})

				By("selecting an option by text", func() {
					selection := page.Find("#some_select")
					Expect(selection.All("option").At(1)).NotTo(BeSelected())
					Expect(selection.Select("second option")).To(Succeed())
					Expect(selection.All("option").At(1)).To(BeSelected())
				})

				By("submitting a form", func() {
					Expect(page.Find("#some_form").Submit()).To(Succeed())
					Eventually(func() bool { return submitted }).Should(BeTrue())
				})
			})
		})

		Describe("Page interactions", func() {
			It("should support retrieving page properties", func() {
				Expect(page).To(HaveTitle("Page Title"))
				Expect(page).To(HaveURL(server.URL + "/"))
				Expect(page.HTML()).To(ContainSubstring("<h1>Title</h1>"))
			})

			It("should support JavaScript", func() {
				By("waiting for page JavaScript to take effect", func() {
					Expect(page.Find("#some_element")).NotTo(HaveText("some text"))
					Eventually(page.Find("#some_element"), "4s").Should(HaveText("some text"))
					Consistently(page.Find("#some_element")).Should(HaveText("some text"))
				})

				// NOTE: disabled due to recent Firefox regression with passing args
				if name != "Firefox" {
					By("executing provided JavaScript", func() {
						arguments := map[string]interface{}{"elementID": "some_element"}
						var result string
						Expect(page.RunScript("return document.getElementById(elementID).innerHTML;", arguments, &result)).To(Succeed())
						Expect(result).To(Equal("some text"))
					})
				}
			})

			It("should support taking screenshots", func() {
				Expect(page.Screenshot(".test.screenshot.png")).To(Succeed())
				defer os.Remove(".test.screenshot.png")
				file, _ := os.Open(".test.screenshot.png")
				_, err := png.Decode(file)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should support links and navigation", func() {
				By("clicking on a link", func() {
					Expect(page.FindByLink("Click Me").Click()).To(Succeed())
					Expect(page.URL()).To(ContainSubstring("#new_page"))
				})

				By("navigating through browser history", func() {
					Expect(page.Back()).To(Succeed())
					Expect(page.URL()).NotTo(ContainSubstring("#new_page"))
					Expect(page.Forward()).To(Succeed())
					Expect(page.URL()).To(ContainSubstring("#new_page"))
				})

				By("refreshing the page", func() {
					checkbox := page.Find("#some_checkbox")
					Expect(checkbox.Check()).To(Succeed())
					Expect(page.Refresh()).To(Succeed())
					Expect(checkbox).NotTo(BeSelected())
				})
			})

			// NOTE: browsers besides PhantomJS do not support JavaScript logs
			if name == "PhantomJS" {
				It("should support retrieving logs", func() {
					Eventually(page).Should(HaveLoggedInfo("some log"))
					Expect(page).NotTo(HaveLoggedError())
					Eventually(page, "4s").Should(HaveLoggedError("ReferenceError: Can't find variable: doesNotExist\n  (anonymous function)"))
				})
			}

			It("should support switching frames", func() {
				By("switching to an iframe", func() {
					Expect(page.Find("#frame").SwitchToFrame()).To(Succeed())
					Expect(page.Find("body")).To(MatchText("Example Domain"))
				})

				// NOTE: PhantomJS does not support Page.SwitchToParentFrame
				if name != "PhantomJS" {
					By("switching back to the default frame by referring to the parent frame", func() {
						Expect(page.SwitchToParentFrame()).To(Succeed())
						Expect(page.Find("body")).NotTo(MatchText("Example Domain"))
					})

					Expect(page.Find("#frame").SwitchToFrame()).To(Succeed())
				}

				By("switching back to the default frame by referring to the root frame", func() {
					Expect(page.SwitchToRootFrame()).To(Succeed())
					Expect(page.Find("body")).NotTo(MatchText("Example Domain"))
				})
			})

			It("should support switching windows", func() {
				Expect(page.Find("#new_window").Click()).To(Succeed())
				Expect(page).To(HaveWindowCount(2))

				By("switching windows", func() {
					Expect(page.SwitchToWindow("new window")).To(Succeed())
					Expect(page.Find("header")).NotTo(BeFound())
					Expect(page.NextWindow()).To(Succeed())
					Expect(page.Find("header")).To(BeFound())
				})

				By("closing windows", func() {
					Expect(page.CloseWindow()).To(Succeed())
					Expect(page).To(HaveWindowCount(1))
				})
			})

			// NOTE: PhantomJS does not support popup boxes
			if name != "PhantomJS" {
				It("should support popup boxes", func() {
					By("interacting with alert popups", func() {
						Expect(page.Find("#popup_alert").Click()).To(Succeed())
						Expect(page).To(HavePopupText("some alert"))
						Expect(page.ConfirmPopup()).To(Succeed())
					})

					By("interacting with confirm boxes", func() {
						var confirmed bool

						Expect(page.Find("#popup_confirm").Click()).To(Succeed())

						Expect(page.ConfirmPopup()).To(Succeed())
						Expect(page.RunScript("return confirmed;", nil, &confirmed)).To(Succeed())
						Expect(confirmed).To(BeTrue())

						Expect(page.Find("#popup_confirm").Click()).To(Succeed())

						Expect(page.CancelPopup()).To(Succeed())
						Expect(page.RunScript("return confirmed;", nil, &confirmed)).To(Succeed())
						Expect(confirmed).To(BeFalse())
					})

					By("interacting with prompt boxes", func() {
						var promptText string

						Expect(page.Find("#popup_prompt").Click()).To(Succeed())

						Expect(page.EnterPopupText("banana")).To(Succeed())
						Expect(page.ConfirmPopup()).To(Succeed())
						Expect(page.RunScript("return promptText;", nil, &promptText)).To(Succeed())
						Expect(promptText).To(Equal("banana"))
					})
				})
			}

			It("should support manipulating and retrieving cookies", func() {
				Expect(page.SetCookie(&http.Cookie{Name: "webdriver-test-cookie", Value: "webdriver value"})).To(Succeed())
				cookies, err := page.GetCookies()
				Expect(err).NotTo(HaveOccurred())
				cookieNames := []string{cookies[0].Name, cookies[1].Name}
				Expect(cookieNames).To(ConsistOf("webdriver-test-cookie", "browser-test-cookie"))
				Expect(page.DeleteCookie("browser-test-cookie")).To(Succeed())
				Expect(page.GetCookies()).To(HaveLen(1))
				Expect(page.ClearCookies()).To(Succeed())
				Expect(page.GetCookies()).To(HaveLen(0))
			})
		})
	})
}
