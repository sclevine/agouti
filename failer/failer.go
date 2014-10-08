package failer

import "github.com/onsi/ginkgo"

type Failer struct {
	async      bool
	callerSkip int
}

func (f *Failer) Fail(message string) {
	f.Skip()

	if f.async {
		panic(message)
	} else {
		ginkgo.Fail(message, f.callerSkip)
	}
}

func (f *Failer) Skip() {
	f.callerSkip += 1
}

func (f *Failer) UnSkip() {
	f.callerSkip -= 1
}

func (f *Failer) Async() {
	f.async = true
}

func (f *Failer) Sync() {
	f.async = false
}
