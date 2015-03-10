package appium

import "github.com/sclevine/agouti/internal/target"

// MultiSelection is a Selection that can be indexes. See `agouti.MultiSelection`.
type MultiSelection struct {
	Selection
	newSelection func(target.Selectors) *Selection
}

// At finds an element at the provided index. See `agouti.MultiSelection:At`
func (s *MultiSelection) At(index int) *Selection {
	return s.newSelection(target.Selectors(s.Selectors()).At(index))
}
