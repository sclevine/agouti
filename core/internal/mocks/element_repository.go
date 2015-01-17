package mocks

import (
	"github.com/sclevine/agouti/core/internal/selection"
)

type ElementRepository struct {
	GetCall struct {
		Selectors      []selection.Selector
		ReturnElements []selection.Element
		Err            error
	}

	GetExactlyOneCall struct {
		Selectors     []selection.Selector
		ReturnElement selection.Element
		Err           error
	}

	GetAtLeastOneCall struct {
		Selectors      []selection.Selector
		ReturnElements []selection.Element
		Err            error
	}
}

func (e *ElementRepository) Get(selectors []selection.Selector) ([]selection.Element, error) {
	e.GetCall.Selectors = selectors
	return e.GetCall.ReturnElements, e.GetCall.Err
}

func (e *ElementRepository) GetExactlyOne(selectors []selection.Selector) (selection.Element, error) {
	e.GetExactlyOneCall.Selectors = selectors
	return e.GetExactlyOneCall.ReturnElement, e.GetExactlyOneCall.Err
}

func (e *ElementRepository) GetAtLeastOne(selectors []selection.Selector) ([]selection.Element, error) {
	e.GetAtLeastOneCall.Selectors = selectors
	return e.GetAtLeastOneCall.ReturnElements, e.GetAtLeastOneCall.Err
}
