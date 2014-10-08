package failer

type Failer struct {
	FailTest   func(message string, callerSkip ...int)
	async      bool
	callerSkip int
}

func (f *Failer) Fail(message string) {
	f.Down()

	if f.async {
		f.Down()
		panic(message)
	} else {
		callerSkip := f.callerSkip
		f.callerSkip = 0
		f.FailTest(message, callerSkip)
	}
}

func (f *Failer) Down() bool {
	f.callerSkip += 1
	return true
}

func (f *Failer) Up(ignored ...bool) {
	f.callerSkip -= 1
}

func (f *Failer) Async() {
	f.async = true
}

func (f *Failer) Sync() {
	f.async = false
}

func (f *Failer) Reset() {
	f.callerSkip = 0
}
