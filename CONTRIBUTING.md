Contributing
============

Rules for pull requests:

1. Everything within reason must have BDD-style tests.
2. Follow all existing patterns in the codebase.
3. Before issuing a pull-request, please make sure to rebase your branch against master.
   Pull requests will not be merged without this.
4. After issuing your pull request, check Travis CI to make sure that all tests still pass.

Development Setup
-----------------

* Clone the repository
* Follow the README instructions to install Ginkgo, Gomega, PhantomJS, ChromeDriver, and Selenium
* Run all of the tests using: `ginkgo -r .`
* Start developing!
