package mobile

import "github.com/sclevine/agouti/api"

type Element struct {
	api.Element
	Session *Session
}

func (e *Element) SetThing() (string, error) {
	var text string
	if err := e.Send("thing", "GET", nil, &text); err != nil {
		return "", err
	}
	return text, nil
}

// override methods like this and return mobile.Element!
func (e *Element) GetElement(selector Selector) (*Element, error) {
	apiElement, err := e.Element.GetElement(selector)
	if err != nil {
		return nil, err
	}

	return &Element{apiElement, e.Session}, nil
}
