# Contributing

## Pull Requests:

1. Everything (within reason) must have BDD-style tests.
2. Test driving is very strongly encourage.
3. Follow all existing patterns and conventions in the codebase.
4. Before issuing a pull-request, please rebase your branch against master.
   If you are okay with the maintainer rebasing your pull request, please mention this.
5. After issuing your pull request, check Travis CI to make sure that all tests still pass.

## Development Setup

* Clone the repository
* Follow the README.md instructions to install Ginkgo, Gomega, PhantomJS, ChromeDriver, and Selenium
* Run all of the tests using: `ginkgo -r .`
* Start developing!

## Method Naming Conventions

### Page or Selection Level

* `Name` - Methods that retrieve data or perform some action should not start with "Get", "Is", or "Set".
* `SetName` - Methods that set data and have a corresponding `Name` method should start with "Set".

### API Level

All API method names should be as close to their endpoint names as possible.
* `GetName` for all GET requests returning a non-boolean
* `IsName` for all GET requests returning a boolean
* `SetName` for POST requests that have matching GET requests
* `Name` for POST requests that perform some action or retrieve data
* `Get<type>Element` for all POST requests returning an element
