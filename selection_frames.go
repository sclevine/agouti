package agouti

import (
	"fmt"

	"github.com/sclevine/agouti/api"
)

// SwitchToFrame focuses on the frame specified by the selection. All new and
// existing selections will refer to the new frame. All further Page methods
// will apply to this frame as well.
func (s *Selection) SwitchToFrame() error {
	selectedElement, err := s.elements.GetExactlyOne(s.selectors)
	if err != nil {
		return fmt.Errorf("failed to select '%s': %s", s, err)
	}

	if err := s.session.Frame(selectedElement.(*api.Element)); err != nil {
		return fmt.Errorf("failed to switch to frame '%s': %s", s, err)
	}
	return nil
}
