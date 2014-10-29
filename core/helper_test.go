package core

import "github.com/sclevine/agouti/core/internal/types"

func TestingPage(client types.Client) Page {
	return &page{client}
}

func TestingSelection(client types.Client) Selection {
	return &selection{client: client}
}

func TestingMultiSelection(client types.Client) MultiSelection {
	return &multiSelection{&selection{client: client}}
}

func TestingDriver(service types.Service) WebDriver {
	return &driver{service: service}
}
