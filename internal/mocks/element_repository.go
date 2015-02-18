package mocks

import (
	"github.com/sclevine/agouti/internal/element"
	"github.com/sclevine/agouti/internal/target"
)

type ElementRepository struct {
	GetCall struct {
		Selectors      target.Selectors
		ReturnElements []element.Element
		Err            error
	}

	GetExactlyOneCall struct {
		Selectors     target.Selectors
		ReturnElement element.Element
		Err           error
	}

	GetAtLeastOneCall struct {
		Selectors      target.Selectors
		ReturnElements []element.Element
		Err            error
	}
}

func (e *ElementRepository) Get(selectors target.Selectors) ([]element.Element, error) {
	e.GetCall.Selectors = selectors
	return e.GetCall.ReturnElements, e.GetCall.Err
}

func (e *ElementRepository) GetExactlyOne(selectors target.Selectors) (element.Element, error) {
	e.GetExactlyOneCall.Selectors = selectors
	return e.GetExactlyOneCall.ReturnElement, e.GetExactlyOneCall.Err
}

func (e *ElementRepository) GetAtLeastOne(selectors target.Selectors) ([]element.Element, error) {
	e.GetAtLeastOneCall.Selectors = selectors
	return e.GetAtLeastOneCall.ReturnElements, e.GetAtLeastOneCall.Err
}
