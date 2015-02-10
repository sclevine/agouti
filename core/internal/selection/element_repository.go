package selection

import (
	"errors"
	"fmt"

	"github.com/sclevine/agouti/api"
)

type ElementRepository struct {
	Client ElementClient
}

type ElementClient interface {
	GetElement(selector api.Selector) (*api.Element, error)
	GetElements(selector api.Selector) ([]*api.Element, error)
}

type Element interface {
	ElementClient
	GetText() (string, error)
	GetAttribute(attribute string) (string, error)
	GetCSS(property string) (string, error)
	IsSelected() (bool, error)
	IsDisplayed() (bool, error)
	IsEnabled() (bool, error)
	IsEqualTo(other *api.Element) (bool, error)
	Click() error
	Clear() error
	Value(text string) error
	Submit() error
}

func (e *ElementRepository) GetAtLeastOne(selectors []Selector) ([]Element, error) {
	elements, err := e.Get(selectors)
	if err != nil {
		return nil, err
	}

	if len(elements) == 0 {
		return nil, errors.New("no elements found")
	}

	return elements, nil
}

func (e *ElementRepository) GetExactlyOne(selectors []Selector) (Element, error) {
	elements, err := e.GetAtLeastOne(selectors)
	if err != nil {
		return nil, err
	}

	if len(elements) > 1 {
		return nil, fmt.Errorf("method does not support multiple elements (%d)", len(elements))
	}

	return elements[0], nil
}

func (e *ElementRepository) Get(selectors []Selector) ([]Element, error) {
	if len(selectors) == 0 {
		return nil, errors.New("empty selection")
	}

	lastElements, err := retrieveElements(e.Client, selectors[0])
	if err != nil {
		return nil, err
	}

	for _, selector := range selectors[1:] {
		elements := []Element{}
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

func retrieveElements(client ElementClient, selector Selector) ([]Element, error) {
	if selector.Single {
		elements, err := client.GetElements(selector.API())
		if err != nil {
			return nil, err
		}

		if len(elements) == 0 {
			return nil, errors.New("element not found")
		} else if len(elements) > 1 {
			return nil, errors.New("ambiguous find")
		}

		return []Element{Element(elements[0])}, nil
	}

	if selector.Indexed && selector.Index > 0 {
		elements, err := client.GetElements(selector.API())
		if err != nil {
			return nil, err
		}

		if selector.Index >= len(elements) {
			return nil, errors.New("element index out of range")
		}

		return []Element{Element(elements[selector.Index])}, nil
	}

	if selector.Indexed && selector.Index == 0 {
		element, err := client.GetElement(selector.API())
		if err != nil {
			return nil, err
		}
		return []Element{Element(element)}, nil
	}

	elements, err := client.GetElements(selector.API())
	if err != nil {
		return nil, err
	}

	newElements := []Element{}
	for _, element := range elements {
		newElements = append(newElements, element)
	}

	return newElements, nil
}
