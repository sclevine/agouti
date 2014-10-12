package dsl

import (
	"github.com/onsi/ginkgo"
	"github.com/sclevine/agouti/page"
)

type Page interface {
	page.PageOnly
	page.Selection
}

func CreatePage() Page {
	page, err := page.PhantomPage()
	if err != nil {
		ginkgo.Fail(err.Error(), 1)
	}

	return page
}
