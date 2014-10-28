package selection

import (
	"errors"
	"fmt"

	"github.com/sclevine/agouti/core/internal/types"
)

type retriever interface {
	GetElements(selector types.Selector) ([]types.Element, error)
	GetElement(selector types.Selector) (types.Element, error)
}

func (s *Selection) getSelectedElements() ([]types.Element, error) {
	elements, err := s.getElements()
	if err != nil {
		return nil, err
	}

	if len(elements) == 0 {
		return nil, errors.New("no elements found")
	}

	return elements, nil
}

func (s *Selection) getSelectedElement() (types.Element, error) {
	elements, err := s.getSelectedElements()
	if err != nil {
		return nil, err
	}

	if len(elements) > 1 {
		return nil, fmt.Errorf("method does not support multiple elements (%d)", len(elements))
	}

	return elements[0], nil
}

func (s *Selection) getElements() ([]types.Element, error) {
	if len(s.selectors) == 0 {
		return nil, errors.New("empty selection")
	}

	lastElements, err := retrieveElements(s.Client, s.selectors[0])
	if err != nil {
		return nil, err
	}

	for _, selector := range s.selectors[1:] {
		elements := []types.Element{}
		for _, element := range lastElements {
			subElements, err := retrieveElements(element, selector)
			if err != nil {
				return nil, err
			}

			elements = append(elements, subElements...)
		}
		lastElements = elements
	}
	return lastElements, nil
}

func retrieveElements(retriever retriever, selector types.Selector) ([]types.Element, error) {
	if selector.Single {
		elements, err := retriever.GetElements(selector)
		if err != nil {
			return nil, err
		}

		if len(elements) == 0 {
			return nil, errors.New("element not found")
		} else if len(elements) > 1 {
			return nil, errors.New("ambiguous find")
		}

		return []types.Element{elements[0]}, nil
	}

	if selector.Indexed && selector.Index > 0 {
		elements, err := retriever.GetElements(selector)
		if err != nil {
			return nil, err
		}

		if selector.Index >= len(elements) {
			return nil, errors.New("element index out of range")
		}

		return []types.Element{elements[selector.Index]}, nil
	}

	if selector.Indexed && selector.Index == 0 {
		element, err := retriever.GetElement(selector)
		if err != nil {
			return nil, err
		}
		return []types.Element{element}, nil
	}

	return retriever.GetElements(selector)
}
