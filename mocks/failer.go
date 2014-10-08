package mocks

type Failer struct {
	Message     string
	DownCount   int
	UpCount     int
	AsyncCalled bool
	SyncCalled  bool
	IsAsync     bool
	ResetCalled bool
}

func (f *Failer) Fail(message string) {
	f.Message = message
	panic("FAILED")
}

func (f *Failer) Down() {
	f.DownCount += 1
}

func (f *Failer) Up() {
	f.UpCount += 1
}

func (f *Failer) Async() {
	f.IsAsync = true
	f.AsyncCalled = true
}

func (f *Failer) Sync() {
	f.IsAsync = false
	f.SyncCalled = true
}


func (f *Failer) Reset() {
	f.ResetCalled = true
}
