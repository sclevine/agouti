package mocks

type Failer struct {
	Message    string
	CallerSkip int
	Failed     bool
}

func (f *Failer) Fail(message string, callerSkip ...int) {
	f.Failed = true
	f.Message = message
	if len(callerSkip) > 0 {
		f.CallerSkip = callerSkip[0]
	}
}
