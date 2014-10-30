package core

import (
	"github.com/sclevine/agouti/core/internal/selection"
	"github.com/sclevine/agouti/core/internal/types"
)

func TestingPage(client types.Client) Page {
	return &page{&baseSelection{&selection.Selection{Client: client}}}
}

func TestingDriver(service types.Service) WebDriver {
	return &driver{service: service}
}
