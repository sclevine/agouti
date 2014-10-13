package dsl

import (
	"github.com/onsi/ginkgo"
	"github.com/sclevine/agouti/page"
)

type AgoutiPage interface {
	page.PageOnly
	page.Selection
}

func CreatePage(pageType ...string) AgoutiPage {
	var (
		newPage AgoutiPage
		err     error
	)

	if len(pageType) == 0 {
		pageType = append(pageType, "phantom")
	}

	switch pageType[0] {
	case "phantom", "phantomjs":
		newPage, err = page.PhantomPage()
	case "chrome":
		newPage, err = page.ChromePage()
	default:
		newPage, err = page.SeleniumPage(pageType[0])
	}

	if err != nil {
		ginkgo.Fail(err.Error(), 1)
	}

	return newPage
}
