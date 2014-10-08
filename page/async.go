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
	a.selection.failer.Down()
	a.executeDeferred(func() {
		a.selection.failer.Down()
		a.selection.ContainText(text)
	})
}

func (a *async) HaveAttribute(attribute, value string) {
	a.selection.failer.Down()
	a.executeDeferred(func() {
		a.selection.failer.Down()
		a.selection.HaveAttribute(attribute, value)
	})
}

func (a *async) executeDeferred(matcher func()) {
	a.selection.failer.Async()
	a.selection.failer.Down()

	timeoutChan := time.After(a.timeout)
	defer a.retry(timeoutChan, matcher)
	matcher()
	a.selection.failer.Sync()
	a.selection.failer.Reset()
}

func (a *async) retry(timeoutChan <-chan time.Time, matcher func()) {
	a.selection.failer.Down()

	if failure := recover(); failure != nil {
		select {
		case <-timeoutChan:
			a.selection.failer.Sync()
			a.selection.failer.Fail(fmt.Sprintf("After %s:\n %s", a.timeout, failure))
		default:
			time.Sleep(a.interval)
			defer a.retry(timeoutChan, matcher)
			matcher()
		}
	}
	a.selection.failer.Sync()
	a.selection.failer.Reset()
}
