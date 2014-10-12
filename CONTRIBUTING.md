Contributing
============

Agouti welcomes all pull-request and feature suggestions as long as they follow a few simple rules:

1. Everything you add must be 100% tested, if it's hard to test ask for help.
2. Follow the patterns that are already laid out in the codebase.
3. Before issuing a pull-request make sure you have rebased against master.  Pull-request will not be merged without this.
4. Watch Travis and make sure all the tests still pass.

Setting Up
----------

* Clone the repo down (keeping the path the same)
* Make sure to grab both ginkgo and gomega (using go get)
* Run all the tests and make sure they pass
* Start developing!

How to Test
-----------

Jump into the top-level directory and run

```go
ginkgo -r .
```