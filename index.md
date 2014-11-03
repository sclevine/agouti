---
layout: default
title: Agouti
---
[Agouti](http://github.com/sclevine/agouti) is an acceptance testing framework and general-purpose WebDriver API. It is best complemented by the [Ginkgo](http://onsi.github.io/ginkgo/) BDD testing framework and [Gomega](http://onsi.github.io/gomega/) matcher library, but it is designed to be both testing framework and matcher library agnostic.

Most of this document is written with the assumption that you will be using Agouti for acceptance testing with both Ginkgo and Gomega. If you are unfamiliar with these libraries, we recommend consulting their documentation first. See [here](http://onsi.github.io/ginkgo/) and [here](http://onsi.github.io/gomega/).

The [`dsl`](http://godoc.org/github.com/sclevine/agouti/dsl) package is not used in any examples outside of the DSL section. While it does reduce the amount of test setup and provide a [Capybara](http://jnicklas.github.io/capybara/)-like DSL, it exports many more identifiers and results in tests that look quite different from standard Ginkgo tests.

---

##Getting Agouti

Just `go get` it:

    $ go get github.com/sclevine/agouti

Or grab only the packages you plan to use:

    $ go get github.com/sclevine/agouti/core
    $ go get github.com/sclevine/agouti/matchers
    $ go get github.com/sclevine/agouti/dsl

If you plan to write acceptance tests using Ginkgo, or you want to use the [`dsl`](https://godoc.org/github.com/sclevine/agouti/dsl) package:

    $ go get github.com/onsi/ginkgo/ginkgo

If you plan to use the [`matchers`](https://godoc.org/github.com/sclevine/agouti/matchers) package, get Gomega:

    $ go get github.com/onsi/gomega

Next, install any WebDrivers you plan to use. For OS X (using [Homebrew](http://brew.sh)):

    $ brew install phantomjs
    $ brew install chromedriver
    $ brew install selenium-server-standalone

---

##Getting Started: Writing Your First Acceptance Test
Agouti is best used with Ginkgo and Gomega. We'll start by setting up your Go package (`potato`) to work with them. For more information, check out the [Ginkgo docs](http://onsi.github.io/ginkgo/) and [Gomega docs](http://onsi.github.io/ginkgo/).

###Setting up Ginkgo to Run Agouti Tests

    $ cd path/to/potato
    $ ginkgo bootstrap --agouti

This will generate a file named `potato_suite_test.go` containing:

    package potato_test

    import (
        . "github.com/onsi/ginkgo"
        . "github.com/onsi/gomega"
        . "github.com/sclevine/agouti/core"

        "testing"
    )

    func TestPotato(t *testing.T) {
        RegisterFailHandler(Fail)
        RunSpecs(t, "Potato Suite")
    }

    var driver WebDriver

    var _ = BeforeSuite(func() {
        var err error

        // Choose a WebDriver:

        driver, err = PhantomJS()
        // driver, err = Selenium()
        // driver, err = Chrome()

        Expect(err).NotTo(HaveOccurred())
        Expect(driver.Start()).To(Succeed())
    })

    var _ = AfterSuite(func() {
        driver.Stop()
    })

Update this file to run your choice of WebDriver. For this example, we'll use Selenium:

    package potato_test

    import (
        . "github.com/onsi/ginkgo"
        . "github.com/onsi/gomega"
        . "github.com/sclevine/agouti/core"

        "testing"
    )

    func TestPotato(t *testing.T) {
        RegisterFailHandler(Fail)
        RunSpecs(t, "Potato Suite")
    }

    var driver WebDriver

    var _ = BeforeSuite(func() {
        var err error
        driver, err = Selenium()
        Expect(err).NotTo(HaveOccurred())
        Expect(driver.Start()).To(Succeed())
    })

    var _ = AfterSuite(func() {
        driver.Stop()
    })

Note that while this setup does not need to be in the `*_suite_test.go` file, we strongly recommend that the `WebDriver` be stopped in an `AfterSuite` block so that extra WebDriver processes will not remain running if Ginkgo is unceremoniously terminated. Ginkgo guarantees that the `AfterSuite` block will run before it exits. (If you prefer not to use a global variable for the WebDriver, or if you would like to reduce this setup, check out the [`dsl`](https://godoc.org/github.com/sclevine/agouti/dsl) package.)

At this point you can run your suite without any tests.

    $ ginkgo #or go test
    Running Suite: Potato Suite
    ===========================
    Random Seed: 1378936983
    Will run 0 of 0 specs


    Ran 0 of 0 Specs in 0.000 seconds
    SUCCESS! -- 0 Passed | 0 Failed | 0 Pending | 0 Skipped PASS

    Ginkgo ran 1 suite in 1.309896055s
    Test Suite Passed

###Adding Acceptance Tests

Let's write an acceptance test covering user login:

    $ ginkgo generate --agouti user_login

This will generate a file named `user_login_test.go` containing:

    package potato_test

    import (
        . "/path/to/potato"

        . "github.com/onsi/ginkgo"
        . "github.com/onsi/gomega"
        . "github.com/sclevine/agouti/core"
    )

    var _ = Describe("UserLogin", func() {
        var page Page

        BeforeEach(func() {
            var err error
            page, err = driver.Page()
            Expect(err).NotTo(HaveOccurred())
        })

        AfterEach(func() {
            page.Destroy()
        })
    })

Now let's start your app and tell Agouti to navigate to it. Agouti can test any service that runs in a web browser, but let's assume that `potato` exports `StartMyApp(port int)`, which starts your app on the provided port. We'll also tell Agouti to use Firefox for these tests.

    package potato_test

    import (
        . "/path/to/potato"

        . "github.com/onsi/ginkgo"
        . "github.com/onsi/gomega"
        . "github.com/sclevine/agouti/core"
    )

    var _ = Describe("UserLogin", func() {
        var page Page

        BeforeEach(func() {
            StartMyApp(3000)

            var err error
            page, err = driver.Page(Use().Browser("firefox"))
            Expect(err).NotTo(HaveOccurred())
        })

        AfterEach(func() {
            page.Destroy()
        })

        It("allows a user to log in and log out", func() {
            By("redirecting the user to the login form from the home page", func() {
                Expect(page.Navigate("http://localhost:3000")).To(Succeed())
                Expect(page).To(HaveURL("http://localhost:3000/login"))
            })

            By("allowing the user to fill out the login form and submit it", func() {
                Eventually(page.FindByLabel("E-mail").Fill("spud@example.com")).Should(Succeed())
                Expect(page.FindByLabel("Password").Fill("secret-password")).To(Succeed())
                Expect(page.Find("#remember_me").Check()).To(Succeed())
                Expect(page.Find("#login_form").Submit()).To(Succeed())
            })

            By("directing the user to the dashboard", func() {
                Eventually(page).Should(HaveTitle("Dashboard"))
            })

            By("allowing the sure to view their profile", func() {
                Expect(page.FindByLink("Profile Page").Click()).To(Succeed())
                profile := page.Find("section.profile")
                Eventually(profile.Find(".greeting")).Should(HaveText("Hello Spud!"))
                Expect(profile.Find("img#profile_pic")).To(BeVisible())
            })

            By("allowing the user to log out", func() {
                Expect(page.Find("#logout").Click()).To(Succeed())
                Expect(page).To(HavePopupText("Are you sure?"))
                Expect(page.ConfirmPopup()).To(Succeed())
                Eventually(page).Should(HaveTitle("Login"))
            })
        })
    })

###Notes

- A new `Selection` can be created from an existing `Selection` or from the page using the `Selectable` methods defined [here](http://godoc.org/github.com/sclevine/agouti/core#Selectable).
- The `Selection` interface is very rich. It supports selecting and asserting on one or more elements by CSS selector, XPath, label, and/or link text. A selection may combine any number of selector types.
- The Agouti matchers (like `HaveTitle` and `BeVisible`) rely only on public `Page` and `Selection` methods (lke `Title` and `Visible`).
- Gomega's [asynchronous assertions](http://onsi.github.io/gomega/#making-asynchronous-assertions) such as `Eventually` may be used to wait for the page to load. This is especially useful for testing JavaScript-heavy web applications.
- The [`dsl`](https://godoc.org/github.com/sclevine/agouti/dsl) package contains actions that will immediately fail the running test if they fail. This can reduce the number of `Expect(...).To(Succeed())` assertions. It also re-declares some of the Ginkgo blocks. For instance, `It` becomes `Scenario`.
- As your test suite grows it's liable to get slow.  Ginkgo makes it easy to parallelize your test suite by spreading different `It`s across multiple test processes, and Agouti supports this.  We can make the above example support parallel tests by spinning up the app-server in the `BeforeEach` on a port unique to the Ginkgo node that is running the test: `StartMyApp(3000+GinkgoParallelNode())`.  Now you can run the tests in parallel with `ginkgo -p`.  For more on test parallelization see the [Ginkgo docs on the topic](http://onsi.github.io/ginkgo/#parallel-specs).

##Reference

Agouti is fully documented using GoDoc. See [`core`](http://godoc.org/github.com/sclevine/agouti/core), [`matchers`](http://godoc.org/github.com/sclevine/agouti/matchers), and [`dsl`](http://godoc.org/github.com/sclevine/agouti/dsl).

More extensive documentation (with more examples!) is coming soon.

##External WebDrivers and Sauce Labs Support

Agouti supports connecting to any WebDriver that supports the [WebDriver Wire Protocol](https://code.google.com/p/selenium/wiki/JsonWireProtocol). This can be accomplished using `Connect` in [`core`](http://godoc.org/github.com/sclevine/agouti/core):

    page := Connect(Use().Browser("safari"), "http://example.com:1234/wd/hub")
    ...
    page.Destroy() // end session

For easy [Sauce Labs](http://saucelabs.com) support, use `SauceLabs`. Note that this does not currently support Sauce Connect.

    page := SauceLabs("my test", "Linux", "firefox", "33", "my-username", "secret-api-key")
    ...
    page.Destroy() // end session

##The Agouti DSL

Agouti provides a [`dsl`](http://godoc.org/github.com/sclevine/agouti/dsl) package that reduces test setup and provides user interactions that fail tests automatically if they are unsuccessful. It is similar to the [Capybara](https://github.com/jnicklas/capybara) DSL, and intended for those who are familiar with Capybara or Cucumber.

Note that the [`dsl`](http://godoc.org/github.com/sclevine/agouti/dsl) package exports numerous identifiers that may conflict with your package identifiers (if both are dot-imported). It also only supports a single running WebDriver process (PhantomJS, Selenium, OR ChromeDriver). If you use Ginkgo for unit testing, your Agouti tests will look like very different things using the [`dsl`](http://godoc.org/github.com/sclevine/agouti/dsl) package. For these reasons, we do not necessarily encourage its use.

That said, you may re-write the above login test using the [`dsl`](http://godoc.org/github.com/sclevine/agouti/dsl) package like so:

`potato_suite_test.go`:

    package potato_test

    import (
        . "github.com/onsi/ginkgo"
        . "github.com/onsi/gomega"
        . "github.com/sclevine/agouti/dsl"

        "testing"
    )

    func TestPotato(t *testing.T) {
        RegisterFailHandler(Fail)
        RunSpecs(t, "Potato Suite")
    }

    var _ = BeforeSuite(func() {
        StartSelenium()
    })

    var _ = AfterSuite(func() {
        StopWebdriver()
    })

`user_login_test.go`:

    package potato_test

    import (
        . "/path/to/potato"

        . "github.com/onsi/ginkgo"
        . "github.com/onsi/gomega"
        . "github.com/sclevine/agouti/core"
        . "github.com/sclevine/agouti/dsl"
    )

    var _ = Feature("UserLogin", func() {
        var page Page

        Background(func() {
            StartMyApp(3000)
            page = CreatePage(Use().Browser("firefox"))
        })

        AfterEach(func() {
            page.Destroy()
        })

        Scenario("allows a user to log in and log out", func() {
            Step("redirecting the user to the login form from the home page", func() {
                Expect(page.Navigate("http://localhost:3000")).To(Succeed())
                Expect(page).To(HaveURL("http://localhost:3000/login"))
            })

            Step("allowing the user to fill out the login form and submit it", func() {
                Fill(page.FindByLabel("E-mail"), "spud@example.com")
                Fill(page.FindByLabel("Password"), "secret-password")
                Check(page.Find("#remember_me"))
                Submit(page.Find("#login_form"))
            })

            Step("directing the user to the dashboard", func() {
                Eventually(page).Should(HaveTitle("Dashboard"))
            })

            Step("allowing the user to view their profile", func() {
                Click(page.FindByLink("Profile Page"))
                profile := page.Find("section.profile")
                Eventually(profile.Find(".greeting")).Should(HaveText("Hello Spud!"))
                Expect(profile.Find("img#profile_pic")).To(BeVisible())
            })

            Step("allowing the user to log out", func() {
                Click(page.Find("#logout"))
                Expect(page).To(HavePopupText("Are you sure?"))
                Expect(page.ConfirmPopup()).To(Succeed())
                Eventually(page).Should(HaveTitle("Login"))
            })
        })
    })

Like Ginkgo test blocks, [`dsl`](http://godoc.org/github.com/sclevine/agouti/dsl) test blocks may be pended by prepending an `X` or `P` or focused by prepending an `F`. For example, an `XScenario` will never run, but a `FFeature` will prevent specs outside of it from running unless they are also focused.

##Using Agouti Without Ginkgo or Gomega

The [`core`](http://godoc.org/github.com/sclevine/agouti/core) package does not depend on Ginkgo or Gomega, and can even be used as a general-purpose WebDriver API for uses outside of acceptance testing.

Here is a partial login test that does not depend on Ginkgo or Gomega. To make this even more Go-like, we'll import the core package as `agouti` instead of dot-importing it.

    package potato_test

    import (
        . "/path/to/potato"
        agouti "github.com/sclevine/agouti/core"
        "testing"
    )

    func TestUserLoginPrompt(t *testing.T) {
        driver := agouti.Selenium()
        page := driver.Page(agouti.Use().Browser("firefox"))

        StartMyApp(3000)

        if err := page.Navigate("http://localhost:3000"); err != nil {
            t.Error("Failed to navigate!")
        }

        loginURL, err := page.URL()
        if err != nil {
            t.Error("Failed to get page URL!")
        }

        if loginURL != "http://localhost:3000/login" {
            t.Error("Expected URL to be http://localhost:3000/login, got ", loginURL)
        }

        loginPrompt, err := page.Find("#prompt").Text()
        if err != nil {
            t.Error("Failed to get login prompt text!")
        }

        if loginPrompt != "Please login!" {
            t.Error("Expected login prompt to be 'Please login!', got ", loginPrompt)
        }

        driver.Stop() // destroys page automatically
    }

If you'd like to use the Agouti matchers without using Ginkgo, you can use Gomega's `RegisterTestingT()` function to use Gomega matchers in Go's XUnit style tests.  See Gomega's [docs for more details](http://onsi.github.io/gomega/#using-gomega-with-golangs-xunit-style-tests).
