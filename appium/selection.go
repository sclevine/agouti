package appium

import (
	"github.com/sclevine/agouti/internal/element"
	"github.com/sclevine/agouti/internal/target"
)

type elementRepository interface {
	Get(selectors target.Selectors) ([]element.Element, error)
	GetAtLeastOne(selectors target.Selectors) ([]element.Element, error)
	GetExactlyOne(selectors target.Selectors) (element.Element, error)
}
