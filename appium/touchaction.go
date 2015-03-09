package appium

import (
	"fmt"
	"strings"

	"github.com/sclevine/agouti/api"
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

type ActionOptions struct {
	// TODO: check which means what, what are the differences between ms and duration ?
	Duration    int    `json:"duration,omitempty"` // which units ??
	Millisecond int    `json:"ms,omitempty"`       // duplicates with Duration ??
	X           int    `json:"x,omitempty"`
	Y           int    `json:"y,omitempty"`
	Element     string `json:"element,omitempty"` // element ID
	Count       int    `json:"count,omitempty"`   // meaning ??
}

func (ma *TouchAction) Tap() *TouchAction {
	return ma.appendAction(newAction("tap", ma.Element))
}

func (ma *TouchAction) Press() *TouchAction {
	return ma
}

func (ma *TouchAction) Release() *TouchAction {
	return ma
}

func (ma *TouchAction) Wait() *TouchAction {
	return ma
}

func (ma *TouchAction) MoveTo() *TouchAction {
	return ma
}

func (ma *TouchAction) Perform() error {
	var actions []interface{}
	for _, act := range ma.Actions {
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
		actions = append(actions, act.Action)
	}
	return strings.Join(actions, ", ")
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
