package dsl

import (
	"fmt"

	"github.com/onsi/ginkgo"
)

type AgoutiFailHandler func(message string, callerSkip ...int)

var globalFailHandler AgoutiFailHandler

func init() {
	globalFailHandler = ginkgo.Fail
}

// RegisterAgoutiFailHandler connects the implied assertions in Agouti's dsl with
// Gingko. When set to ginkgo.Fail (the default), failures in Agouti's dsl-provided
// methods will cause test failures in Ginkgo.
func RegisterAgoutiFailHandler(handler AgoutiFailHandler) {
	globalFailHandler = handler
}

func checkFailure(err error) {
	if err != nil {
		globalFailHandler(fmt.Sprintf("Agouti failure: %s", err), 2)
	}
}
