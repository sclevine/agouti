package mobile

import "github.com/sclevine/agouti/api"

type Element struct {
	*api.Element
	Session *Session
}

func (e *Element) SetThing() (string, error) {
	var text string
	if err := e.Send("thing", "GET", nil, &text); err != nil {
		return "", err
	}
	return text, nil
}

func (e *Element) Tap() error {
	ma := newTouchAction(e.Session)
	ma.Element = e.Element
	return ma.Perform()
}

func (e *Element) Actions() *TouchAction {
	ma := newTouchAction(e.Session)
	ma.Element = e.Element
	return ma
}

func (e *Element) PerformMultiTouch(actions ...*TouchAction) error {
	return nil
}

// override methods like this and return mobile.Element!
func (e *Element) GetElement(selector api.Selector) (*Element, error) {
	apiElement, err := e.Element.GetElement(selector)
	if err != nil {
		return nil, err
	}

	return &Element{apiElement, e.Session}, nil
}
