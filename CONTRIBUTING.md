Contributing
============

Rules for pull requests:

1. Everything within reason must have BDD-style tests.
2. Follow all existing patterns in the codebase.
3. Before issuing a pull-request, please make sure to rebased your branch against master.
   Pull requests will not be merged without this.
4. After issuing your pull request, look at Travis CI to make sure all tests still pass.

Setting Up
----------

* Clone the repository
* Follow the README instructions to install Ginkgo, Gomega, PhantomJS, and Selenium
* Run all the tests using: `ginkgo -r .`
* Start developing!