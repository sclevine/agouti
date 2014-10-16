package dsl

import (
	"github.com/onsi/ginkgo"
	"github.com/sclevine/agouti/core"
	"fmt"
)

func Click(selection core.Selection) {
	check(selection.Click())
}

func DoubleClick(selection core.Selection) {
	check(selection.DoubleClick())
}

func Fill(selection core.Selection, text string) {
	check(selection.Fill(text))
}

func Check(selection core.Selection) {
	check(selection.Check())
}

func Uncheck(selection core.Selection) {
	check(selection.Uncheck())
}

func Select(selection core.Selection, text string) {
	check(selection.Select(text))
}

func Submit(selection core.Selection) {
	check(selection.Submit())
}

func check(err error) {
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("Action failed: %s", err))
	}
}
