package mocks

type Failer struct {
	Message    string
	CallerSkip int
}

func (f *Failer) Fail(message string, callerSkip ...int) {
	f.Message = message
	if len(callerSkip) > 0 {
		f.CallerSkip = callerSkip[0]
	}
	panic("FAILED")
}
