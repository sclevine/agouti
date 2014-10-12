package matchers

import (
	"github.com/onsi/gomega/types"
	"github.com/sclevine/agouti/matchers/internal/page"
)

func HaveTitle(title string) types.GomegaMatcher {
	return &page.HaveTitleMatcher{ExpectedTitle: title}
}
