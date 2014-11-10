package dsl

import (
	"fmt"

	"github.com/onsi/ginkgo"
	"github.com/sclevine/agouti/core"
)

// Click is comparable to Expect(selection.Click()).To(Succeed())
func Click(selection core.Selection) {
	check(selection.Click())
}

// DoubleClick is comparable to Expect(selection.DoubleClick()).To(Succeed())
func DoubleClick(selection core.Selection) {
	check(selection.DoubleClick())
}

// Fill is comparable to Expect(selection.Fill(text)).To(Succeed())
func Fill(selection core.Selection, text string) {
	check(selection.Fill(text))
}

// Check is comparable to Expect(selection.Check()).To(Succeed())
func Check(selection core.Selection) {
	check(selection.Check())
}

// Uncheck is comparable to Expect(selection.Uncheck()).To(Succeed())
func Uncheck(selection core.Selection) {
	check(selection.Uncheck())
}

// Select is comparable to Expect(selection.Select(text)).To(Succeed())
func Select(selection core.Selection, text string) {
	check(selection.Select(text))
}

// Submit is comparable to Expect(selection.Submit()).To(Succeed())
func Submit(selection core.Selection) {
	check(selection.Submit())
}

// SwitchToFrame is comparable to Expect(selection.SwitchToFrame()).To(Succeed())
func SwitchToFrame(selection core.Selection) {
	check(selection.SwitchToFrame())
}

// SwitchToParentFrame is comparable to Expect(page.SwitchToParentFrame()).To(Succeed())
func SwitchToParentFrame(page core.Page) {
	check(page.SwitchToParentFrame())
}

// SwitchToRootFrame is comparable to Expect(page.SwitchToRootFrame()).To(Succeed())
func SwitchToRootFrame(page core.Page) {
	check(page.SwitchToRootFrame())
}

func check(err error) {
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("Action failed: %s", err), 2)
	}
}
