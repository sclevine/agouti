package page

import (
	"fmt"
	"time"
)

type async struct {
	selection *selection
	timeout   time.Duration
	interval  time.Duration
}

func (a *async) Selector() string {
	return a.selection.Selector()
}

func (a *async) ContainText(text string) {
	a.selection.page.failer.Async()
	a.selection.page.failer.Skip()

	timeoutChan := time.After(a.timeout)
	matcher := func() {
		a.selection.page.failer.Skip()
		a.selection.ContainText(text)
	}
	defer a.retry(timeoutChan, matcher)
	matcher()
}

func (a *async) retry(timeoutChan <-chan time.Time, matcher func()) {
	a.selection.page.failer.Skip()
	if failure := recover(); failure != nil {
		select {
		case <-timeoutChan:
			a.selection.page.failer.Sync()
			a.selection.page.failer.UnSkip()
			a.selection.page.failer.Fail(fmt.Sprintf("After 0.5 seconds:\n %s", failure))
		default:
			time.Sleep(a.interval)
			defer a.retry(timeoutChan, matcher)
			matcher()
		}
	}
}
