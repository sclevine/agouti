package mobile

import (
	"github.com/sclevine/agouti/api"
)

type Session struct {
	*api.Session
}

func (s *Session) SetEndpoint(thing string) error {
	if err := s.Send("endpoint", "POST", thing, nil); err != nil {
		return err
	}

	return nil
}

// override methods like this and return mobile.Element!
func (s *Session) GetElement(selector Selector) (*Element, error) {
	apiElement, err := s.Session.GetElement(selector)
	if err != nil {
		return nil, err
	}

	return &Element{apiElement, s}, nil
}
