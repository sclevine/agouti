package dsl

import "github.com/sclevine/agouti/core"

// SwitchToFrame is comparable to Expect(selection.SwitchToFrame()).To(Succeed())
func SwitchToFrame(selection core.Selection) {
	checkFailure(selection.SwitchToFrame())
}

// Click is comparable to Expect(selection.Click()).To(Succeed())
func Click(selection core.Selection) {
	checkFailure(selection.Click())
}

// DoubleClick is comparable to Expect(selection.DoubleClick()).To(Succeed())
func DoubleClick(selection core.Selection) {
	checkFailure(selection.DoubleClick())
}

// Fill is comparable to Expect(selection.Fill(text)).To(Succeed())
func Fill(selection core.Selection, text string) {
	checkFailure(selection.Fill(text))
}

// Check is comparable to Expect(selection.Check()).To(Succeed())
func Check(selection core.Selection) {
	checkFailure(selection.Check())
}

// Uncheck is comparable to Expect(selection.Uncheck()).To(Succeed())
func Uncheck(selection core.Selection) {
	checkFailure(selection.Uncheck())
}

// Select is comparable to Expect(selection.Select(text)).To(Succeed())
func Select(selection core.Selection, text string) {
	checkFailure(selection.Select(text))
}

// Submit is comparable to Expect(selection.Submit()).To(Succeed())
func Submit(selection core.Selection) {
	checkFailure(selection.Submit())
}
