package appium

import (
	"fmt"
	"strings"

	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/api/mobile"
	"github.com/sclevine/agouti/internal/element"
	"github.com/sclevine/agouti/internal/target"
)

type TouchAction struct {
	actions  []action
	elements elementRepository
	session  mobileSession
}

type action struct {
	mobile.Action
	selectors target.Selectors
}

func NewTouchAction(session mobileSession) *TouchAction {
	return &TouchAction{
		elements: element.Repository{Client: session},
		session:  session,
	}
}

func (a *action) String() string {
	out := []string{}
	opts := a.Options

	if opts.ElementSelection != nil {
		out = append(out, fmt.Sprintf("element='%s'", opts.ElementSelection.String()))
	}
	if opts.X != 0 {
		out = append(out, fmt.Sprintf("x=%d", opts.X))
	}
	if opts.Y != 0 {
		out = append(out, fmt.Sprintf("y=%d", opts.Y))
	}
	if opts.Millisecond != 0 {
		out = append(out, fmt.Sprintf("ms=%d", opts.Millisecond))
	}
	if opts.Count != 0 {
		out = append(out, fmt.Sprintf("count=%d", opts.Count))
	}
	if opts.Duration != 0 {
		out = append(out, fmt.Sprintf("duration=%d", opts.Duration))
	}

	return fmt.Sprintf("%s(%s)", a.Action, strings.Join(out, ", "))
}

func (t *TouchAction) Tap() *TouchAction {
	action := mobile.Action{Action: "tap"}
	return t.append(action, nil)
}

func (t *TouchAction) append(action mobile.Action, selectors agouti.Selectors) *TouchAction {
	newAction := action{action, selectors.(target.Selectors)}
	touchAction := NewTouchAction(t.session)
	touchAction.actions = append(t.actions, newAction)
	return touchAction
}

func (t *TouchAction) PressPosition(x, y int) *TouchAction {
	action := mobile.Action{
		Action:  "press",
		Options: {X: x, Y: y},
	}
	return t.append(action, nil)
}

func (t *TouchAction) PressElement(selection *agouti.Selection) *TouchAction {
	action := mobile.Action{Action: "press"}
	return t.append(action, selection.Selectors())
}

func (t *TouchAction) LongPressPosition(x, y, duration int) *TouchAction {
	action := mobile.Action{
		Action:  "longPress",
		Options: {X: x, Y: y, Duration: duration},
	}
	return t.append(action, nil)
}

func (t *TouchAction) LongPressElement(selection *agouti.Selection, duration int) *TouchAction {
	action := mobile.Action{
		Action:  "longPress",
		Options: {Duration: duration},
	}
	return t.append(action, selection.Selectors())
}

func (t *TouchAction) Release() *TouchAction {
	action := mobile.Action{Action: "release"}
	return t.append(action, nil)
}

func (t *TouchAction) Wait(ms int) *TouchAction {
	action := mobile.Action{
		Action:  "longPress",
		Options: {Millisecond: ms},
	}
	return t.append(action, nil)
}

func (t *TouchAction) MoveToPosition(x, y int) *TouchAction {
	action := mobile.Action{
		Action:  "moveTo",
		Options: {X: x, Y: y},
	}
	return t.append(action, nil)
}

func (t *TouchAction) MoveToElement(selection *agouti.Selection) *TouchAction {
	action := mobile.Action{Action: "moveTo"}
	return t.append(action, selection.Selectors())
}

func (t *TouchAction) Perform() error {
	var actions []mobile.Action

	for _, action := range t.actions {

		// resolve elements if present
		if action.selectors != nil {
			selectedElement, err := t.elements.GetExactlyOne(action.selectors)
			if err != nil {
				return fmt.Errorf("failed to retrieve element for selection '%s': %s", action.selectors, err)
			}
			action.Options.Element = selectedElement.(*api.Element).ID
		}

		actions = append(actions, action.Action)
	}

	if err := t.session.PerformTouch(actions); err != nil {
		return fmt.Errorf("error performing touch actions '%s': %s", t, err)
	}
	return nil
}

func (ma *TouchAction) String() string {
	var actions []string
	for _, act := range ma.Actions {
		actions = append(actions, act.String())
	}
	return strings.Join(actions, " -> ")
}
