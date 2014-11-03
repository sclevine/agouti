# Contributing

## Rules for pull requests:

1. Everything within reason must have BDD-style tests.
2. Test driving is very strongly encourage.
2. Follow all existing patterns and coventions in the codebase.
3. Before issuing a pull-request, please make sure to rebase your branch against master.
   If you are okay with the maintainer rebasing your pull, please mention this in the request.
4. After issuing your pull request, check Travis CI to make sure that all tests still pass.

## Development Setup

* Clone the repository
* Follow the README instructions to install Ginkgo, Gomega, PhantomJS, ChromeDriver, and Selenium
* Run all of the tests using: `ginkgo -r .`
* Start developing!

## Method Naming Conventions

### Page Level

* `Name` - Methods that retrieve data or perform some action should not start with "Get", "Is", or "Set".
* `SetName` - Methods that set data and have a corresponding `Name` method should start with "Set".

### API level

All API method names should be as close to their endpoint name as possible.
* `GetName` for all GET requests returning a non-boolean
* `IsName` for all GET requests returning a boolean
* `SetName` for POST requests that have matching GET requests
* `Name` for POST requests that perform some action or retrieve data
* `Get<type>Element` for all POST requests returning an element
