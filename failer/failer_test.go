package failer_test

import (
	. "github.com/sclevine/agouti/failer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Failer", func() {
	var (
		failer *Failer
		failMessage string
		failCallerSkip int
	)


	BeforeEach(func() {
		failer = &Failer{FailTest: func(message string, callerSkip ...int) {
			failMessage = message
			failCallerSkip = callerSkip[0]
		}}
	})

	Describe("#Fail", func() {
		Context("when asynchronous", func() {
			BeforeEach(func() {
				failer.Async()
			})

			It("panics with a provided message", func() {
				defer func() {
					Expect(recover()).To(Equal("async panic"))
				}()
				failer.Fail("async panic")
			})

			It("increments the caller skip by two", func() {
				Expect(func() { failer.Fail("async panic") }).To(Panic())
				failer.Sync()
				failer.Fail("async panic")
				Expect(failCallerSkip).To(Equal(3))
			})
		})

		Context("when not asynchronous", func() {
			It("fails the test", func() {
				failer.Fail("test failure")
				Expect(failMessage).To(Equal("test failure"))
			})

			It("increments the caller count by one, then clears it", func() {
				failer.Fail("test failure")
				Expect(failCallerSkip).To(Equal(1))
				failer.Fail("test failure")
				Expect(failCallerSkip).To(Equal(1))
			})
		})
	})

	Describe("#Down", func() {
		It("increments the caller skip", func() {
			failer.Down()
			failer.Fail("test failure")
			Expect(failCallerSkip).To(Equal(2))
		})
	})

	Describe("#Up", func() {
		It("decrements the caller skip", func() {
			failer.Up()
			failer.Fail("test failure")
			Expect(failCallerSkip).To(Equal(0))
		})
	})

	Describe("#Sync/#Async", func() {
		It("turns off #Async panicing when #Sync is called", func() {
			failer.Async()
			failer.Sync()
			failer.Fail("test failure")
			Expect(failMessage).To(Equal("test failure"))
		})
	})
})
