package mocks

type Failer struct {
	Message    string
	CallerSkip int
	IsAsync      bool
}

func (f *Failer) Fail(message string) {
	f.Message = message

	panic("FAILED")
}

func (f *Failer) Skip() {
	f.CallerSkip += 1
}

func (f *Failer) UnSkip() {
	f.CallerSkip -= 1
}

func (f *Failer) Async() {
	f.IsAsync = true
}

func (f *Failer) Sync() {
	f.IsAsync = false
}
