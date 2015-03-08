---
layout: default
title: Agouti
---
[Agouti](http://github.com/sclevine/agouti) is an universal WebDriver client for Go. For acceptance or integration testing, it is best complemented by the [Ginkgo](http://onsi.github.io/ginkgo/) BDD testing framework and [Gomega](http://onsi.github.io/gomega/) matcher library, but it is designed to be both testing-framework- and matcher-library-agnostic.

Much of this document is written with the assumption that you will be using Agouti for acceptance testing with both Ginkgo and Gomega. If you are unfamiliar with these libraries, consult their documentation first. See [here](http://onsi.github.io/ginkgo/) and [here](http://onsi.github.io/gomega/). Note that the [`agouti`](https://godoc.org/github.com/sclevine/agouti) package can be used by itself as a general-purpose WebDriver client for Go.

---

##Getting Agouti

Just `go get` it:

    $ go get github.com/sclevine/agouti

If you plan to write acceptance tests using Ginkgo:

    $ go get github.com/onsi/ginkgo/ginkgo

If you plan to use the Gomega matchers provided by the [`matchers`](https://godoc.org/github.com/sclevine/agouti/matchers) package, get Gomega:

    $ go get github.com/onsi/gomega

Next, install any WebDrivers you plan to use. For Mac OS X (using [Homebrew](http://brew.sh)):

    $ brew install phantomjs
    $ brew install chromedriver
    $ brew install selenium-server-standalone
    
(Consider running `brew update` before these commands.)

We currently support PhantomJS 1.9.7+, Selenium WebDriver 2.44.0+, and ChromeDriver 2.13+. See [this thread](https://code.google.com/p/selenium/issues/detail?can=2&q=7933&colspec=ID%20Stars%20Type%20Status%20Priority%20Milestone%20Owner%20Summary&id=7933) if you have issues running Selenium with Safari on Mac OS X. Any WebDriver conforming to the [WebDriver Wire Protocol](https://code.google.com/p/selenium/wiki/JsonWireProtocol) should (theoretically) work with Agouti, and can be configured using [`agouti.NewWebDriver`](https://godoc.org/github.com/sclevine/agouti#NewWebDriver).

---

##Getting Started: Writing Your First Acceptance Test
For acceptance or integration testing, Agouti is best used with Ginkgo and Gomega. We'll start by setting up your Go package (named `potato`) to work with them. (For more information, check out the [Ginkgo docs](http://onsi.github.io/ginkgo/) and [Gomega docs](http://onsi.github.io/ginkgo/)).

###Bootstrapping Ginkgo to Run Agouti Tests

    $ cd path/to/potato
    $ ginkgo bootstrap --agouti

This will generate a file named `potato_suite_test.go` containing:

    package potato_test

    import (
        "testing"
        
        . "github.com/onsi/ginkgo"
        . "github.com/onsi/gomega"
        "github.com/sclevine/agouti"
    )

    func TestPotato(t *testing.T) {
        RegisterFailHandler(Fail)
        RunSpecs(t, "Potato Suite")
    }

    var agoutiDriver *agouti.WebDriver

    var _ = BeforeSuite(func() {
        // Choose a WebDriver:

        agoutiDriver = agouti.PhantomJS()
        // agoutiDriver = agouti.Selenium()
        // agoutiDriver = agouti.ChromeDriver()

        Expect(agoutiDriver.Start()).To(Succeed())
    })

    var _ = AfterSuite(func() {
        Expect(agoutiDriver.Stop()).To(Succeed())
    })

Update this file with your choice of WebDriver. For this example, we'll use Selenium.

    package potato_test

    import (
        "testing"
        
        . "github.com/onsi/ginkgo"
        . "github.com/onsi/gomega"
        "github.com/sclevine/agouti"
    )

    func TestPotato(t *testing.T) {
        RegisterFailHandler(Fail)
        RunSpecs(t, "Potato Suite")
    }

    var agoutiDriver *agouti.WebDriver

    var _ = BeforeSuite(func() {
        agoutiDriver = agouti.Selenium()
        Expect(agoutiDriver.Start()).To(Succeed())
    })

    var _ = AfterSuite(func() {
        Expect(agoutiDriver.Stop()).To(Succeed())
    })

Note that while this setup does not need to be in the `potato_suite_test.go` file, we strongly recommend that the `*agouti.WebDriver` be stopped in an `AfterSuite` block so that extra WebDriver processes will not remain running if Ginkgo is unceremoniously terminated. Ginkgo guarantees that the `AfterSuite` block will run before it exits.

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

Let's write an acceptance test covering user login. First, use Ginkgo to generate an Agouti test template:

    $ ginkgo generate --agouti user_login

This will generate a file named `user_login_test.go` containing:

    package potato_test

    import (
        . "path/to/potato"
        . "github.com/onsi/ginkgo"
        . "github.com/onsi/gomega"
        . "github.com/sclevine/agouti/matchers"
        "github.com/sclevine/agouti"
    )

    var _ = Describe("UserLogin", func() {
        var page *agouti.Page

        BeforeEach(func() {
            var err error
            page, err = agoutiDriver.NewPage()
            Expect(err).NotTo(HaveOccurred())
        })

        AfterEach(func() {
            Expect(page.Destroy()).To(Succeed())
        })
    })

Now let's start your application and tell Agouti to navigate to it. Agouti can test any service that runs in a web browser, but let's assume that `potato` exports `StartMyApp(port int)`, which starts your application on the provided port. We'll tell Agouti to use Firefox for these tests.

    package potato_test

    import (
        . "path/to/potato"
        . "github.com/onsi/ginkgo"
        . "github.com/onsi/gomega"
        . "github.com/sclevine/agouti/matchers"
        "github.com/sclevine/agouti"
    )

    var _ = Describe("UserLogin", func() {
        var page *agouti.Page

        BeforeEach(func() {
            StartMyApp(3000)
            
            var err error
            page, err = agoutiDriver.NewPage(agouti.Browser("firefox"))
            Expect(err).NotTo(HaveOccurred())
        })

        AfterEach(func() {
            Expect(page.Destroy()).To(Succeed())
        })

        It("should manage user authentication", func() {
            By("redirecting the user to the login form from the home page", func() {
                Expect(page.Navigate("http://localhost:3000")).To(Succeed())
                Expect(page).To(HaveURL("http://localhost:3000/login"))
            })

            By("allowing the user to fill out the login form and submit it", func() {
                Eventually(page.FindyByLabel("E-mail")).Should(BeFound())
                Expect(page.FindByLabel("E-mail").Fill("spud@example.com")).To(Succeed())
                Expect(page.FindByLabel("Password").Fill("secret-password")).To(Succeed())
                Expect(page.Find("#remember_me").Check()).To(Succeed())
                Expect(page.Find("#login_form").Submit()).To(Succeed())
            })
            
            By("allowing the user to view their profile", func() {
                Eventually(page.FindByLink("Profile Page")).Should(BeFound())
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

- A new `*agouti.Selection` can be created from a `*agouti.Page`, an existing `*agouti.Selection`, or an existing `*agouti.MultiSelection` using the `Find-`, `First-`, and `All-` methods defined on either type.
- The `*agouti.Selection` methods support selecting and asserting on one or more elements by CSS selector, XPath, label, button text, and/or link text. A selection may combine any number of these selector types.
- The Agouti matchers (like `HaveTitle` and `BeVisible`) rely only on public `*agouti.Page` and `*agouti.Selection` methods (lke `Title` and `Visible`).
- Gomega's [asynchronous assertions](http://onsi.github.io/gomega/#making-asynchronous-assertions) such as `Eventually` may be used to wait for the page to load. This is especially useful for testing JavaScript-heavy web applications.
- As your test suite grows, it's liable to get slow. Ginkgo makes it easy to parallelize your test suite by spreading different `It` blocks across multiple test processes, and Agouti supports this. We can make the above example support parallel tests by spinning up the app-server in the `BeforeEach` on a port unique to the Ginkgo node that is running the test: `StartMyApp(3000+GinkgoParallelNode())`. After adjusting the corresponding URLs in the `It` block, you can run the tests in parallel with `ginkgo -p`. For more on test parallelization see the [Ginkgo docs on the topic](http://onsi.github.io/ginkgo/#parallel-specs).
- While using many `It` blocks that each create and destroy a `*agouti.Page` is an effective way to parallelize your tests, it can actually slow your test suite down if you don't need this parallelization. Calling `agoutiDriver.NewPage()` and `page.Destroy()` in your `BeforeSuite` and `AfterSuite` blocks (or using a few large `It` blocks with many `By` blocks) can eliminate the overhead of creating a new WebDriver session for each test.

##Reference

Agouti is fully documented using GoDoc. See [`agouti`](http://godoc.org/github.com/sclevine/agouti) and [`matchers`](http://godoc.org/github.com/sclevine/agouti/matchers).

The [`api`](http://godoc.org/github.com/sclevine/agouti/api) package provides low-level access to the WebDriver. It currently does not have a fixed API, but this will change in the near future (with the addition of adequate documentation).

More extensive documentation (with more examples!) coming soon.

##External WebDrivers and Sauce Labs Support

Agouti supports managing any WebDriver that supports the [WebDriver Wire Protocol](https://code.google.com/p/selenium/wiki/JsonWireProtocol) and that is launched by a command running a foreground process. This can be complished using `agouti.NewWebDriver`:

    command := []string{"java", "-jar", "selenium-server.jar", "-port", "{{.Port}}"}
    driver := NewWebDriver("http://{{.Address}}/wd/hub", command)
    Expect(driver.Start()).To(Succeed())
    page, err := driver.NewPage()
    ...
    Expect(page.Destroy()).To(Succeed()) // end session

Agouti supports connecting to any running WebDriver that supports the [WebDriver Wire Protocol](https://code.google.com/p/selenium/wiki/JsonWireProtocol). This can be accomplished using `Connect` in [`core`](http://godoc.org/github.com/sclevine/agouti/core):

    page, err := agouti.NewPage("http://example.com:1234/wd/hub")
    ...
    Expect(page.Destroy()).To(Succeed()) // end session

For easy [Sauce Labs](http://saucelabs.com) support, use `SauceLabs`. Note that this does not currently support Sauce Connect.

    page, err := SauceLabs("my test", "Linux", "firefox", "33", "my-username", "secret-api-key")
    ...
    Expect(page.Destroy()).To(Succeed()) // end session

##Using Agouti with Gomega and XUnit Tests

If you would prefer to use Go's built-in XUnit tests instead of Ginkgo, the [`agouti`](http://godoc.org/github.com/sclevine/agouti) and [`matchers`](http://godoc.org/github.com/sclevine/agouti/matchers) packages make this easy.

To use Agouti with Gomega and XUnit style tests, check out this simple example:

    package potato_test

    import (
        "testing"
        
        . "path/to/potato"
        . "github.com/onsi/gomega"
        . "github.com/sclevine/agouti/matchers"
        "github.com/sclevine/agouti"
    )

    func TestUserLoginPrompt(t *testing.T) {
        RegisterTestingT(t)

        driver := agouti.Selenium()
        Expect(driver.Start()).To(Succeed())
        page, err := driver.NewPage(agouti.Browser("firefox"))
        Expect(err).NotTo(HaveOccurred())

        StartMyApp(3000)

        Expect(page.Navigate("http://localhost:3000")).To(Succeed())
        Expect(page).To(HaveURL("http://localhost:3000"))
        Expect(page.Find("#prompt")).To(HaveText("Please login!"))

        Expect(driver.Stop()).To(Succeed()) // calls page.Destroy() automatically
    }

See Gomega's [docs for more details](http://onsi.github.io/gomega/#using-gomega-with-golangs-xunit-style-tests). Note that using Agouti without Ginkgo will not allow you to run your specs in parallel.

###Without Dot-Imports

This is the most Go-like way of using Agouti for acceptance testing.

    package potato_test

    import (
        "testing"
        
        "path/to/potato"
        "github.com/sclevine/agouti"
        am "github.com/sclevine/agouti/matchers"
        gm "github.com/onsi/gomega"
    )

    func TestUserLoginPrompt(t *testing.T) {
        gm.RegisterTestingT(t)

        driver := agouti.Selenium()
        gm.Expect(driver.Start()).To(gm.Succeed())
        page, err := driver.NewPage(agouti.Browser("firefox"))
        gm.Expect(err).NotTo(gm.HaveOccurred())

        potato.StartMyApp(3000)

        gm.Expect(page.Navigate("http://localhost:3000")).To(gm.Succeed())
        gm.Expect(page).To(am.HaveURL("http://localhost:3000"))
        gm.Expect(page.Find("#prompt")).To(am.HaveText("Please login!"))

        gm.Expect(driver.Stop()).To(gm.Succeed()) // calls page.Destroy() automatically
    }

Alternatively:

    package potato_test

    import (
        "testing"
            
        "path/to/potato"
        "github.com/sclevine/agouti/core"
        "github.com/sclevine/agouti/matchers"
        "github.com/onsi/gomega"
    )

    Expect := gomega.Expect
    Succeed := gomega.Succeed
    HaveOccurred := gomega.HaveOccurred
    HaveText := matchers.HaveText
    HaveURL := matchers.HaveURL

    func TestUserLoginPrompt(t *testing.T) {
        gomega.RegisterTestingT(t)

        driver := agouti.Selenium()
        Expect(driver.Start()).To(Succeed())
        page, err := driver.NewPage(agouti.Browser("firefox"))
        Expect(err).NotTo(HaveOccurred())

        potato.StartMyApp(3000)

        Expect(page.Navigate("http://localhost:3000")).To(Succeed())
        Expect(page).To(HaveURL("http://localhost:3000"))
        Expect(page.Find("#prompt")).To(HaveText("Please login!"))

        Expect(driver.Stop()).To(Succeed()) // calls page.Destroy() automatically
    }

###Using Agouti by Itself

The [`agouti`](http://godoc.org/github.com/sclevine/agouti) package by itself does not depend on Ginkgo or Gomega. It can be used as a general-purpose WebDriver client.

Here is a part of a login test that does not depend on Ginkgo or Gomega.

    package potato_test

    import (
        "testing"

        "path/to/potato"
        "github.com/sclevine/agouti"
    )

    func TestUserLoginPrompt(t *testing.T) {
        driver := agouti.Selenium()
        page, err := driver.NewPage(agouti.Browser("firefox"))
        if err != nil {
            t.Error("Failed to open page.")
        }

        potato.StartMyApp(3000)

        if err := page.Navigate("http://localhost:3000"); err != nil {
            t.Error("Failed to navigate.")
        }

        loginURL, err := page.URL()
        if err != nil {
            t.Error("Failed to get page URL.")
        }

        expectedLoginURL := "http://localhost:3000/login"
        if loginURL != expectedLoginURL {
            t.Error("Expected URL to be", expectedLoginURL, "but got", loginURL)
        }

        loginPrompt, err := page.Find("#prompt").Text()
        if err != nil {
            t.Error("Failed to get login prompt text.")
        }

        expectedPrompt := "Please login."
        if loginPrompt != expectedPrompt {
            t.Error("Expected login prompt to be", expectedPrompt, "but got", loginPrompt)
        }

        if err := driver.Stop(); err != nil {
            t.Error("Failed close open pages and stop WebDriver.")
        }
    }

