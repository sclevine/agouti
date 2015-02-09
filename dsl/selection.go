package dsl

type ActionSelection interface {
	SwitchToFrame() error
	Click() error
	DoubleClick() error
	Fill(text string) error
	Check() error
	Uncheck() error
	Select(text string) error
	Submit() error
}

// SwitchToFrame is comparable to Expect(selection.SwitchToFrame()).To(Succeed())
func SwitchToFrame(selection ActionSelection) {
	checkFailure(selection.SwitchToFrame())
}

// Click is comparable to Expect(selection.Click()).To(Succeed())
func Click(selection ActionSelection) {
	checkFailure(selection.Click())
}

// DoubleClick is comparable to Expect(selection.DoubleClick()).To(Succeed())
func DoubleClick(selection ActionSelection) {
	checkFailure(selection.DoubleClick())
}

// Fill is comparable to Expect(selection.Fill(text)).To(Succeed())
func Fill(selection ActionSelection, text string) {
	checkFailure(selection.Fill(text))
}

// Check is comparable to Expect(selection.Check()).To(Succeed())
func Check(selection ActionSelection) {
	checkFailure(selection.Check())
}

// Uncheck is comparable to Expect(selection.Uncheck()).To(Succeed())
func Uncheck(selection ActionSelection) {
	checkFailure(selection.Uncheck())
}

// Select is comparable to Expect(selection.Select(text)).To(Succeed())
func Select(selection ActionSelection, text string) {
	checkFailure(selection.Select(text))
}

// Submit is comparable to Expect(selection.Submit()).To(Succeed())
func Submit(selection ActionSelection) {
	checkFailure(selection.Submit())
}
