package page

import (
	"time"
	"fmt"
)

type async struct {
	selection *selection
}

func (a *async) ContainText(text string) {
	timeoutChan := time.After(500 * time.Millisecond)
	matcher := func() { a.selection.ContainText(text) }
	defer a.retry(timeoutChan, matcher)
	matcher()
}

func (a *async) retry(timeoutChan <-chan time.Time, matcher func()) {
	if failure := recover(); failure != nil {
		select {
		case <-timeoutChan:
			a.selection.page.fail(fmt.Sprintf("After 0.5 seconds:\n %s", failure), 100)
		default:
			time.Sleep(100 * time.Millisecond)
			defer a.retry(timeoutChan, matcher)
			matcher()
		}
	}
}
