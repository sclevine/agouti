package dsl_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/dsl"
)

var _ = Describe("DSL sanity checks", func() {
	Feature("Background", func() {
		var runsBackground bool

		Background(func() {
			runsBackground = true
		})

		Scenario("Background is a Ginkgo BeforeEach", func() {
			Expect(runsBackground).To(BeTrue())
		})

		Scenario("Step is a Ginkgo By", func() {
			var stepRuns bool

			Step("steps are run", func() {
				Expect(stepRuns).To(BeFalse())
				stepRuns = true
			})

			Expect(stepRuns).To(BeTrue())
		})
	})

	if os.Getenv("DSL_PENDING_SANITY_CHECKS") == "true" {
		XFeature("this Describe is pending (using X)", func() {
			Scenario("so this would not run", func() {
				Fail("failed to pend spec")
			})
		})

		PFeature("this Describe is pending (using P)", func() {
			Scenario("so this would not run", func() {
				Fail("failed to pend spec")
			})
		})

		XScenario("this is pending (using X) and would not run", func() {
			Fail("failed to pend spec")
		})

		PScenario("this is pending (using P) and would not run", func() {
			Fail("failed to pend spec")
		})
	}

	if os.Getenv("DSL_FOCUSED_SANITY_CHECKS") == "true" {
		FFeature("this Describe is focused", func() {
			Scenario("so this will run", func() {
			})
		})

		FScenario("this is focused and will run", func() {
		})

		Scenario("this is not focused and will not run", func() {
			Fail("failed to focus specs")
		})
	}
})
