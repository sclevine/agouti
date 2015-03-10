package appium

import (
	"fmt"
	"strings"

	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/target"
)

type TouchAction struct {
	Actions []Action
	Element *api.Element
	session mobileSession
}

func NewTouchAction(session mobileSession) *TouchAction {
	return &TouchAction{
		Actions: make([]Action, 0),
		session: session,
	}
}

type Action struct {
	Action  string        `json:"action"`
	Options ActionOptions `json:"options,omitempty"`
}

func (a *Action) String() string {
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

type ActionOptions struct {
	// TODO: check which means what, what are the differences between ms and duration ?
	Duration         int        `json:"duration,omitempty"` // which units ??
	Millisecond      int        `json:"ms,omitempty"`       // duplicates with Duration ??
	X                int        `json:"x,omitempty"`
	Y                int        `json:"y,omitempty"`
	Element          string     `json:"element,omitempty"` // element ID
	ElementSelection *Selection `json:"-"`                 // element, to be resolved as an element ID before performing
	Count            int        `json:"count,omitempty"`   // meaning ??
}

func (ma *TouchAction) Tap() *TouchAction {
	return ma.appendAction(newAction("tap", ma.Element))
}

func (ma *TouchAction) PressPosition(x, y int) *TouchAction {
	act := newAction("press", ma.Element)
	act.Options.X = x
	act.Options.Y = y
	return ma.appendAction(act)
}

func (ma *TouchAction) PressElement(sel *Selection) *TouchAction {
	act := newAction("press", ma.Element)
	act.Options.ElementSelection = sel
	return ma.appendAction(act)
}

func (ma *TouchAction) LongPressPosition(x, y, duration int) *TouchAction {
	act := newAction("longPress", ma.Element)
	act.Options.X = x
	act.Options.Y = y
	act.Options.Duration = duration
	return ma.appendAction(act)
}

func (ma *TouchAction) LongPressElement(sel *Selection, duration int) *TouchAction {
	act := newAction("longPress", ma.Element)
	act.Options.ElementSelection = sel
	act.Options.Duration = duration
	return ma.appendAction(act)
}



func (ma *TouchAction) Release() *TouchAction {
	return ma.appendAction(newAction("release", ma.Element))
}

func (ma *TouchAction) Wait(ms int) *TouchAction {
	act := newAction("wait", ma.Element)
	act.Options.Millisecond = ms
	return ma.appendAction(act)
}

func (ma *TouchAction) MoveToPosition(x, y int) *TouchAction {
	act := newAction("moveTo", ma.Element)
	act.Options.X = x
	act.Options.Y = y
	return ma.appendAction(act)
}

func (ma *TouchAction) MoveToElement(sel *Selection) *TouchAction {
	act := newAction("moveTo", ma.Element)
	act.Options.ElementSelection = sel
	return ma.appendAction(act)
}

func (ma *TouchAction) Perform() error {
	var actions []interface{}

	for _, act := range ma.Actions {

		// resolve elements if present
		if act.Options.ElementSelection != nil {
			sel := act.Options.ElementSelection
			el, err := sel.elements.GetExactlyOne(target.Selectors(sel.Selectors()))
			if err != nil {
				return fmt.Errorf("failed to resolve element on action '%s': %s", act.Action, err)
			}
			act.Options.Element = el.GetID()
		}

		actions = append(actions, interface{}(act))
	}

	err := ma.session.PerformTouch(actions)
	if err != nil {
		return fmt.Errorf("error performing touch actions '%s': %s", ma, err)
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

func (ma *TouchAction) appendAction(action Action) *TouchAction {
	ma.Actions = append(ma.Actions, action)
	return ma
}

func newAction(actionType string, element *api.Element) Action {
	action := Action{
		Action:  actionType,
		Options: ActionOptions{},
	}
	if element != nil {
		action.Options.Element = element.ID
	}
	return action
}
